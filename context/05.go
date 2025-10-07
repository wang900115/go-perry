package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dowork(ctx)
}

func dowork(ctx context.Context) {
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped early:", ctx.Err())
			return
		default:
			fmt.Println("Working on step", i)
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("Work complete")
}
