package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			warDeclaration := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}

			err := pubsub.PublishJSON(
			publishCh,
			routing.ExchangePerilTopic,
			routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
			warDeclaration,
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}


func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(warRec gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(warRec)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			err := pubsub.PublishGob(
			publishCh,
			routing.ExchangePerilTopic,
			routing.GameLogSlug+"."+warRec.Attacker.Username,
			fmt.Sprintf("%s won a war against %s", winner, loser),
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			err := pubsub.PublishGob(
				publishCh,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+warRec.Attacker.Username,
				fmt.Sprintf("%s won a war against %s", winner, loser),
				)
				if err != nil {
					return pubsub.NackRequeue
				}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			err := pubsub.PublishGob(
				publishCh,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+warRec.Attacker.Username,
				fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser),
				)
				if err != nil {
					return pubsub.NackRequeue
				}
			return pubsub.Ack
		default:
			fmt.Println("Error in war")
			return pubsub.NackDiscard
		}
	}
}