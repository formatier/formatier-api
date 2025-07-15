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
	mx           sync.Mutex
	amqpConn     *amqp.Connection
	amqpMainChan *amqp.Channel
	consumers    []*Consumer
}

func (conn *Connection) AutoConnect(ctx context.Context, path string, amqpConfig *amqp.Config) {
	if amqpConfig == nil {
		amqpConfig = &amqp.Config{}
	}

	amqpConn, err := amqp.DialConfig(path, *amqpConfig)
	if err != nil {
		panic(fmt.Sprintf("Error on first connection %v", err))
	}
	conn.amqpConn = amqpConn

	conn.amqpMainChan, err = amqpConn.Channel()
	if err != nil {
		log.Println("Error when tring to connect: %w", err)
	}

	reconnectChan := make(chan struct{})

	go func() {
		for {
			select {
			case <-amqpConn.NotifyClose(make(chan *amqp.Error)):
				fmt.Println("Noticed close, waiting for 2 seconds, please wait...")
				if amqpConn.IsClosed() {
					reconnectChan <- struct{}{}
				}
				time.Sleep(2 * time.Second)
			case <-ctx.Done():
				conn.amqpMainChan.Close()
				close(reconnectChan)
				return
			}
		}
	}()

	go func() {
		for range reconnectChan {
			log.Println("Waiting for new connection, please wait...")
			conn.mx.Lock()

			log.Println("Connecting to server...")
			amqpConn, err := amqp.DialConfig(path, *amqpConfig)
			if err != nil {
				log.Printf("Cannot reconnect to server: %s", err)
				amqpConn.Close()
				conn.mx.Unlock()
				continue
			}
			log.Println("Successfully connected to server!")

			log.Println("Creating a new main chanel.")
			amqpMainChan, err := amqpConn.Channel()
			if err != nil {
				log.Printf("Cannot reconnect to server: %s", err)
				amqpConn.Close()
				conn.mx.Unlock()
				continue
			}
			log.Println("Successfully created a new main chanel!")

			log.Println("Reconsuming comsumers...")
			for _, consumer := range conn.consumers {
				err := consumer.Reconnect(amqpConn)
				if err != nil {
					log.Printf("Cannot reconnect to server: %s", err)
					conn.mx.Unlock()
					continue
				}
			}
			log.Println("Successfully recomsumed consumers!")

			conn.amqpConn = amqpConn
			conn.amqpMainChan = amqpMainChan

			fmt.Println("Successfully reconnected")
			conn.mx.Unlock()
		}
	}()
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
