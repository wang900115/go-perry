package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// queue declaration
	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	// Create a new Gin instance
	router := gin.Default()

	// Add route to send a message to the queue

	router.GET("/send", func(c *gin.Context) {
		msg := c.Query("msg")
		if msg == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mesasage is required"})
			return
		}

		// Create a message
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		}
		// Publish a message to queue
		err = channel.Publish("", queueName, false, false, message)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish messsage"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": msg, "status": "success"})
	})

	log.Fatal(router.Run(":8080"))
}
