package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		fmt.Println("準備發送")
		ch <- 42
		fmt.Println("發送完畢")
	}()

	fmt.Println("準備接收")
	num := <-ch
	fmt.Println("收到:", num)
}
