package main

import "fmt"

func main() {
	ch := make(chan string)

	go func() {
		ch <- "Hello from goroutine!"
	}()

	msg := <-ch
	fmt.Println(msg)
}
