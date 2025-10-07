package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "來自ch1 的訊息"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "來自ch2 的訊息"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		case <-time.After(2 * time.Second):
			fmt.Println("超時了")
		}

	}

	fmt.Println("完成")
}
