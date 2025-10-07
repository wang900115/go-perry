package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d 退出: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Worker %d 工作中\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker(ctx, 1)
	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println("主程序結束")
}
