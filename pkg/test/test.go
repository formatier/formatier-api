package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(
		"nats://127.0.0.1:4222",
		nats.CustomReconnectDelay(func(attempts int) time.Duration {
			return time.Duration(attempts+1) * time.Second
		}),
		nats.RetryOnFailedConnect(true),
	)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "UserCommands",
		Replicas:  1,
		Retention: nats.WorkQueuePolicy,
		Subjects:  []string{"user.>"},
	})
	if err != nil {
		panic(err)
	}

	_, err = js.AddConsumer(
		"UserCommands",
		&nats.ConsumerConfig{
			Durable:   "user-command",
			AckPolicy: nats.AckExplicitPolicy,
		},
	)
	if err != nil {
		panic(err)
	}

	sub, err := js.PullSubscribe(
		"user.read.byID",
		"user-command",
	)

	if err != nil {
		panic(err)
	}
	defer sub.Drain()

	for {
		msgBatch, err := sub.FetchBatch(5)
		if err != nil {
			panic(err)
		}

		for msg := range msgBatch.Messages() {
			fmt.Printf("To: %s, Message %s\n", msg.Subject, string(msg.Data))
			time.Sleep(time.Duration(float32(time.Second) * 0.2))

			if msg.Reply == "" {
				fmt.Println("⚠️  No reply subject - this is the problem!")
				err := msg.AckSync()
				if err != nil {
					panic(err)
				}
				continue
			}

			err := msg.Respond([]byte("Hello World"))
			if err != nil {
				panic(err)
			}

			err = msg.AckSync()
			if err != nil {
				panic(err)
			}
		}
	}
}
