package main

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

var enqueuer = work.NewEnqueuer("demo_app", redisPool)

func main() {
	_, err := enqueuer.Enqueue("email",
		work.Q{"userID": 10, "subject": "testing"})

	if err != nil {
		log.Fatal(err)
	}

	enqueuer.Enqueue("report",
		work.Q{"userID": 5})

}
