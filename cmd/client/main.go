package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	const connString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	defer conn.Close()
	fmt.Println("Peril client connected to RabbitMQ")

	// Prompt the user for a username
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("no username given")
	}

	queueName := routing.PauseKey + "." + username
	fmt.Printf("the queue name is: %s\n", queueName)
	ch, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, 1) // 1 for "transient" in the last arg
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}
	// This print statement is a placeholder because I am not using the channel or queue yet in this file
	fmt.Println(ch, queue)

	gs := gamelogic.NewGameState(username)

	// Create REPL with an infinite loop:

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gs.CommandSpawn(words)
			if err != nil{
				log.Fatalf("invalid use of spawn: %v", err)
			}
		case "move":
			_, err := gs.CommandMove(words)
			if err != nil {
				log.Fatalf("invalid use of move: %v", err)
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming is not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command")
		}

	}
}
