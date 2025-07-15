package rabbitmq

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerConfig struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func NewConsumer(
	deliveryMessageChan chan<- *amqp.Delivery,
	amqpConnection *amqp.Connection,
	consumerConfig *ConsumerConfig,
) (*Consumer, error) {
	reconsumeChan := make(chan struct{}, 1)

	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumerConfig: consumerConfig,
		reconsumeChan:  reconsumeChan,

		amqpChannel:         amqpChannel,
		deliveryMessageChan: deliveryMessageChan,
	}, nil
}

type Consumer struct {
	mx sync.Mutex

	consumerConfig *ConsumerConfig
	reconsumeChan  chan struct{}

	amqpChannel         *amqp.Channel
	deliveryMessageChan chan<- *amqp.Delivery
}

func (w *Consumer) Reconnect(newConnection *amqp.Connection) error {
	w.mx.Lock()
	defer w.mx.Unlock()

	if w.amqpChannel != nil {
		w.amqpChannel.Close()
	}

	newChannel, err := newConnection.Channel()
	if err != nil {
		return err
	}
	w.amqpChannel = newChannel

	select {
	case w.reconsumeChan <- struct{}{}:
	default:
	}

	return nil
}

func (w *Consumer) Consume(queueName string) error {
	deliveryChanChan := make(chan (<-chan amqp.Delivery))

	go func() {
		for range w.reconsumeChan {
			deliveryChan, err := w.amqpChannel.Consume(
				queueName,
				w.consumerConfig.Consumer,
				w.consumerConfig.AutoAck,
				w.consumerConfig.Exclusive,
				w.consumerConfig.NoLocal,
				w.consumerConfig.NoWait,
				w.consumerConfig.Args,
			)
			if err != nil {
				log.Error(err)
				continue
			}
			deliveryChanChan <- deliveryChan
		}
	}()

	go func() {
		for deliveryChan := range deliveryChanChan {
			fmt.Println("new deliveryChan")

			for delivery := range deliveryChan {
				fmt.Println("recived new message")
				w.deliveryMessageChan <- &delivery
			}
			fmt.Println("deliveryChan closed, waiting for new channel...")
		}
	}()

	w.reconsumeChan <- struct{}{}

	return nil
}
