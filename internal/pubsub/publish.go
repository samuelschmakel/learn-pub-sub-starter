package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
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

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	fmt.Printf("Publishing to exchange: %s with key %s\n", exchange, key)
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(val)
	if err != nil {
		fmt.Println("Error encoding data")
		return err
	}

	encodedBytes := buffer.Bytes()
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body: encodedBytes,
	})
}