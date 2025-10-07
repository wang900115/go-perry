package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine cancelled:", ctx.Err())
		}
	}(ctx)

	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(500 * time.Millisecond)
}
