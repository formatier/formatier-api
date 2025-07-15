package rabbitmq

import (
	"context"
	"fmt"

	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConnection() *Connection {
	return &Connection{
		consumers: []*Consumer{},
	}
}

type Connection struct {
	mx           sync.RWMutex
	amqpConn     *amqp.Connection
	amqpMainChan *amqp.Channel
	consumers    []*Consumer

	isConnected    bool
	reconnectDelay time.Duration
	maxRetryDelay  time.Duration
}

func (conn *Connection) AutoConnect(ctx context.Context, path string, amqpConfig *amqp.Config) {
	if amqpConfig == nil {
		amqpConfig = &amqp.Config{}
	}

	conn.reconnectDelay = 2 * time.Second
	conn.maxRetryDelay = 30 * time.Second

	err := conn.connect(path, amqpConfig)
	if err != nil {
		panic(fmt.Sprintf("Error on first connection %v", err))
	}

	go conn.monitorConnection(ctx, path, amqpConfig)
}

func (conn *Connection) connect(path string, amqpConfig *amqp.Config) error {
	conn.mx.Lock()
	defer conn.mx.Unlock()

	if conn.amqpConn != nil && !conn.amqpConn.IsClosed() {
		conn.amqpConn.Close()
	}

	amqpConn, err := amqp.DialConfig(path, *amqpConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	amqpMainChan, err := amqpConn.Channel()
	if err != nil {
		amqpConn.Close()
		return fmt.Errorf("failed to create main channel: %w", err)
	}

	conn.amqpConn = amqpConn
	conn.amqpMainChan = amqpMainChan
	conn.isConnected = true

	log.Println("Successfully connected to RabbitMQ")
	return nil
}

func (conn *Connection) monitorConnection(ctx context.Context, path string, amqpConfig *amqp.Config) {
	for {
		conn.mx.RLock()
		currentConn := conn.amqpConn
		conn.mx.RUnlock()

		if currentConn == nil {
			time.Sleep(5 * time.Second)
			continue
		}

		closeNotify := make(chan *amqp.Error, 1)
		currentConn.NotifyClose(closeNotify)

		select {
		case err := <-closeNotify:
			if err != nil {
				log.Printf("Connection closed: %v", err)
			} else {
				log.Println("Connection closed gracefully")
			}

			conn.mx.Lock()
			conn.isConnected = false
			conn.mx.Unlock()

			go conn.handleReconnect(ctx, path, amqpConfig)
			return
		case <-ctx.Done():
			log.Println("Context cancelled, stopping connection monitor")
			conn.mx.Lock()
			if conn.amqpMainChan != nil {
				conn.amqpMainChan.Close()
			}
			if conn.amqpConn != nil {
				conn.amqpConn.Close()
			}
			conn.mx.Unlock()
			return
		}
	}
}

func (conn *Connection) handleReconnect(ctx context.Context, path string, amqpConfig *amqp.Config) {
	retryDelay := conn.reconnectDelay

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		log.Printf("Attempting to reconnect in %v...", retryDelay)
		time.Sleep(retryDelay)

		log.Println("Connecting to server...")
		err := conn.connect(path, amqpConfig)
		if err != nil {
			log.Printf("Failed to reconnect: %v", err)

			retryDelay = min(retryDelay*2, conn.maxRetryDelay)
			continue
		}

		log.Println("Successfully reconnected!")

		conn.mx.RLock()
		consumers := make([]*Consumer, len(conn.consumers))
		copy(consumers, conn.consumers)

		currentConn := conn.amqpConn
		conn.mx.RUnlock()

		log.Println("Reconnecting consumers...")
		for i, consumer := range consumers {
			err := consumer.Reconnect(currentConn)
			if err != nil {
				log.Printf("Failed to reconnect consumer %d: %v", i, err)
				continue
			}
		}
		log.Println("Successfully reconnected all consumers!")

		go conn.monitorConnection(ctx, path, amqpConfig)
		return
	}
}

type ExchangeConfig struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type ExchangeBinder func(key, destination string, config *ExchangeBindingConfig)

type ExchangeBindingConfig struct {
	NoWait bool
	Args   amqp.Table
}

type ExchangeBinding struct {
	ExchangeBinder ExchangeBinder
	Key            string
	Config         *ExchangeBindingConfig
}

func (conn *Connection) registerExchange(exchangeName string, exchangeMode string, config *ExchangeConfig, bindToExchange []*ExchangeBinding) (queueBinder, exchangeBinder ExchangeBinder) {
	if config == nil {
		config = &ExchangeConfig{}
	}
	err := conn.amqpMainChan.ExchangeDeclare(
		exchangeName,
		exchangeMode,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		panic(err)
	}

	for _, exchangeBinding := range bindToExchange {
		exchangeBinding.ExchangeBinder(exchangeBinding.Key, exchangeName, exchangeBinding.Config)
	}

	exchangeBinder = func(key, destination string, config *ExchangeBindingConfig) {
		if config == nil {
			config = &ExchangeBindingConfig{}
		}
		err := conn.amqpMainChan.ExchangeBind(
			destination,
			key,
			exchangeName,
			config.NoWait,
			config.Args,
		)
		if err != nil {
			panic(err)
		}
	}

	queueBinder = func(key, destination string, config *ExchangeBindingConfig) {
		if config == nil {
			config = &ExchangeBindingConfig{}
		}
		err := conn.amqpMainChan.QueueBind(
			destination,
			key,
			exchangeName,
			config.NoWait,
			config.Args,
		)
		if err != nil {
			panic(err)
		}
	}

	return
}

func (conn *Connection) RegisterDirectExchange(name string, config *ExchangeConfig, bindToExchange ...*ExchangeBinding) (queueBinder, exchangeBinder ExchangeBinder) {
	return conn.registerExchange(name, amqp.ExchangeDirect, config, bindToExchange)
}

func (conn *Connection) RegisterFanoutExchange(name string, config *ExchangeConfig, bindToExchange ...*ExchangeBinding) (queueBinder, exchangeBinder ExchangeBinder) {
	return conn.registerExchange(name, amqp.ExchangeFanout, config, bindToExchange)
}

func (conn *Connection) RegisterTopicExchange(name string, config *ExchangeConfig, bindToExchange ...*ExchangeBinding) (queueBinder, exchangeBinder ExchangeBinder) {
	return conn.registerExchange(name, amqp.ExchangeTopic, config, bindToExchange)
}

func (conn *Connection) RegisterHeadersExchange(name string, config *ExchangeConfig, bindToExchange ...*ExchangeBinding) (queueBinder, exchangeBinder ExchangeBinder) {
	return conn.registerExchange(name, amqp.ExchangeHeaders, config, bindToExchange)
}

type QueueConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func (conn *Connection) RegisterQueue(queueName string, config *QueueConfig, bindToExchange ...*ExchangeBinding) string {
	if config == nil {
		config = &QueueConfig{}
	}
	conn.amqpMainChan.QueueDeclare(
		queueName,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)

	for _, exchangeBinding := range bindToExchange {
		exchangeBinding.ExchangeBinder(exchangeBinding.Key, queueName, exchangeBinding.Config)
	}

	return queueName
}

func (conn *Connection) NewConsumer(deliveryMessageChan chan<- *amqp.Delivery, queueName string, config *ConsumerConfig) error {
	if config == nil {
		config = &ConsumerConfig{}
	}
	consumer, err := NewConsumer(deliveryMessageChan, conn.amqpConn, config)
	if err != nil {
		return fmt.Errorf("cannot register consumer: %w", err)
	}

	consumer.Consume(queueName)
	conn.consumers = append(conn.consumers, consumer)

	return nil
}
