package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "ServicelQueue"

func main() {
	// new connection
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// open channel
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()

	// Subscribe to get messages from the queue
	messages, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case message := <-messages:
			log.Printf("Message: %s\n", message.Body)
		case <-sigChan:
			log.Println("Interrupt detected")
			os.Exit(0)
		}
	}
}
