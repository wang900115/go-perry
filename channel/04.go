package main

import "fmt"

func sendOnly(ch chan<- int) {
	ch <- 42
}

func receiveOnly(ch <-chan int) {
	fmt.Println(<-ch)
}

func main() {
	ch := make(chan int)
	go sendOnly(ch)
	receiveOnly(ch)
}
