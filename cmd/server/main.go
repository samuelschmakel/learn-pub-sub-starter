package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ")

	// Create a new connection (amqp.channel)
	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}
	defer publishCh.Close()

	err = pubsub.PublishJSON(
		publishCh, 
		routing.ExchangePerilDirect, 
		routing.PauseKey, 
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		log.Fatalf("could not publish time: %v", err)
	}

	fmt.Println("Pause message sent!")
}
