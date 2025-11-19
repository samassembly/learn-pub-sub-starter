package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"github.com/samassembly/learn-pub-sub-starter/internal/gamelogic"
	"github.com/samassembly/learn-pub-sub-starter/internal/pubsub"
	"github.com/samassembly/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamestate := gamelogic.NewGameState(username)

	Repl: //label for loop break
	for {
		input := gamelogic.GetInput()
		if input == nil {
			continue
		}

		switch input[0] {
			case "spawn":
				err := gamestate.CommandSpawn(input)
				if err != nil {
					fmt.Println("Failed to spawn unit: %v", err)
				}
			case "move":
				_, err := gamestate.CommandMove(input)
				if err != nil {
					fmt.Println("Invalid move: %v", err)
					break
				}
				fmt.Println("Move Successful!")
			case "status":
				gamestate.CommandStatus()
				continue
			case "help":
				gamelogic.PrintClientHelp()
				continue
			case "spam":
				fmt.Println("Spamming not allowed yet!")
				continue
			case "quit":
				gamelogic.PrintQuit()
				break Repl
			default:
				fmt.Println("Unknown command:", input[0])
				continue
		}
	}


	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}