package main

import (
	"fmt"
	"github.com/samassembly/learn-pub-sub-starter/internal/pubsub"
	"github.com/samassembly/learn-pub-sub-starter/internal/routing"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	//declare connection string
	const connectStr = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectStr)
	if err != nil {
		fmt.Errorf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection Successful")

	publishCh, err := conn.Channel() 
	if err != nil {
		fmt.Errorf("Failed to create channel: %v", err)
	}

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
	fmt.Println("Pause message sent!")
}
