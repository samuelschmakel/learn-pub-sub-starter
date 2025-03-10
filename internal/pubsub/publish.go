package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	fmt.Printf("Publishing to exchange: %s with key: %s\n", exchange, key)
	jsonData, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal json: %v", err)
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: jsonData,
	})
}