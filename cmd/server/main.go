package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connString := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	defer conn.Close()
	fmt.Println("connection was successful")

	// Create a new connection (amqp.channel)
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not open new connection: %v", err)
	}
	defer ch.Close()

	// Publish a message to the exchange
	val := routing.PlayingState{
		IsPaused: true,
	}

	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, val)
	if err != nil {
		log.Fatalf("couldn't publish JSON: %v", err)
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println()
	fmt.Println("Shutting down Peril server...")
}
