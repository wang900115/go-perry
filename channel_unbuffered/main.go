package main

import (
	"fmt"
	"time"
)

func main() {
	orders := make(chan string)

	go func() {
		for i := 1; i <= 5; i++ {
			order := fmt.Sprintf("Coffee order #%d", i)
			orders <- order
			fmt.Println("Place:", order)
		}
		close(orders)
	}()

	for order := range orders {
		fmt.Printf("Preparing: %s\n", order)
		time.Sleep(2 * time.Second)
		fmt.Printf("Served: %s \n", order)
	}
}
