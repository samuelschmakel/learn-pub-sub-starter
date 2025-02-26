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
	ch, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.SimpleQueueTransient) // 1 for "transient" in the last arg
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}
	// This print statement is a placeholder because I am not using the channel or queue yet in this file
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+gs.GetUsername(), routing.PauseKey, pubsub.SimpleQueueTransient, handlerPause(gs))
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	// Bind to the army_moves.* routing key
	queueName = routing.ArmyMovesPrefix + "." + username
	err = pubsub.SubscribeJSON[gamelogic.ArmyMove](conn, routing.ExchangePerilTopic, queueName, routing.ArmyMovesPrefix+".*", pubsub.SimpleQueueTransient, handlerMove(gs))

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
				fmt.Println(err)
				continue
			}
		case "move":
			// Execute the move locally
			moveData, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Publish the move to other clients
			routingKey := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, gs.GetUsername())
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routingKey, moveData)
			if err != nil {
				fmt.Println("Failed to publish move:", err)
				continue
			}

			fmt.Println("Move published successfully")

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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

// To do: write handlerMove function
func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(move)
	}
}