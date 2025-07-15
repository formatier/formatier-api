package main

import (
	"context"
	"formatier-api/pkg/rabbitmq"
	"formatier-api/pkg/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx := context.Background()
	conn := rabbitmq.NewConnection()
	conn.AutoConnect(ctx, "amqp://guest:guest@localhost:5672", nil)

	orchCommandQueueBinder, _ := conn.RegisterDirectExchange("orch_command", &rabbitmq.ExchangeConfig{Durable: false})
	userServiceCommandQueue := conn.RegisterQueue(
		"user_service-commands",
		&rabbitmq.QueueConfig{Durable: true},
		&rabbitmq.ExchangeBinding{
			ExchangeBinder: orchCommandQueueBinder,
			Key:            "user.*",
		},
	)

	deliveryChan := make(chan *amqp.Delivery)
	conn.NewConsumer(deliveryChan, userServiceCommandQueue, &rabbitmq.ConsumerConfig{})

	service.NewRouter(context.TODO()).Listen(deliveryChan)
}
