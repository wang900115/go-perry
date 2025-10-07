package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func acquireLock(client *redis.Client, lockKey string, timeout time.Duration) bool {
	ctx := context.Background()
	lockAcquired, err := client.SetNX(ctx, lockKey, "1", timeout).Result()
	if err != nil {
		fmt.Println("Error acquiring lock: ", err)
		return false
	}
	return lockAcquired
}

func releaseLock(client *redis.Client, lockKey string) {
	ctx := context.Background()
	client.Del(ctx, lockKey)
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

	lockKey := "my_lock"
	lockTimeout := 20 * time.Second

	if acquireLock(client, lockKey, lockTimeout) {
		fmt.Println("Lock acquired succeessfully")
		time.Sleep(20 * time.Second)
		fmt.Println("Work done!")

		releaseLock(client, lockKey)
		fmt.Println("Lock released!")
	} else {
		fmt.Println("Failed to acquire lock, Resource is already locked")
	}
}
