package main

import (
	"fmt"
	"log"

	"github.com/samassembly/learn-pub-sub-starter/internal/pubsub"
	"github.com/samassembly/learn-pub-sub-starter/internal/routing"
	"github.com/samassembly/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()
	Loop:  //label for breaking loop
	for { // This creates an infinite loop
		input := gamelogic.GetInput()
		if input == nil {
			continue
		}

		switch input[0] {
			case "pause":
				fmt.Println("Sending Pause Message...")
				err = pubsub.PublishJSON(
					publishCh,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: true,
					},
				)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
			case "resume":
				fmt.Println("Resuming operation...")
				err = pubsub.PublishJSON(
					publishCh,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: false,
					},
				)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
			case "quit":
				fmt.Println("Quitting program...")
				break Loop
			default:
				fmt.Println("Unknown command:", input[0])
			}
	}
}
