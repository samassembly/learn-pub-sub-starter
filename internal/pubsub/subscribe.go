package pubsub

import (
	"fmt"
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T),
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			var msg T
			if err := json.Unmarshal(m.Body, &msg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			handler(msg)
			if err := m.Ack(false); err != nil {
            	log.Printf("Failed to ack message: %v", err)
        	}
		}
	}
	return nil
}