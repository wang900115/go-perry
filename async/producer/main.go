package main

import (
	"asy/tasks"
	"log"

	"github.com/hibiken/asynq"
)

const redis = "127.0.0.1:6379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redis})

	defer client.Close()

	task, err := tasks.NewEmailTask("test@test.com", "Welcome", "Thank you for signing up")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}

	log.Printf("enqueued task: id= %s queue= %s", info.ID, info.Queue)
}
