package main

import (
	"asy/tasks"
	"log"

	"github.com/hibiken/asynq"
)

const redis = "127.0.0.1:6379"

func main() {
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: redis},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		})

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmail, tasks.EmailTaskHandler)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
