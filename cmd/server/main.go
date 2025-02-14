package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
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

	// Create loop

	for {
		slice := gamelogic.GetInput()
		if len(slice) == 0 {
			continue
		}

		if slice[0] == "pause" {
			fmt.Println("Sending a pause message...")
			err = pubsub.PublishJSON(
				publishCh, 
				routing.ExchangePerilDirect, 
				routing.PauseKey, 
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Fatalf("could not publish message: %v", err)
			}
		
			fmt.Println("Pause message sent!")
		} else if slice[0] == "resume" {
			fmt.Println("Sending a resume message...")
			err = pubsub.PublishJSON(
				publishCh, 
				routing.ExchangePerilDirect, 
				routing.PauseKey, 
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Fatalf("could not publish message: %v", err)
			}
		
			fmt.Println("Resume message sent!")
		} else if slice[0] == "quit" {
			fmt.Println("Exiting...")
			break
		} else {
			fmt.Println("Invalid command")
		}
	}

	// End of loop creation
}
