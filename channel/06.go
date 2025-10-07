package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan<- int, wg *sync.WaitGroup) {
	for i := 1; i <= 5; i++ {
		fmt.Println("生產:", i)
		ch <- i
		time.Sleep(1000 * time.Millisecond)
	}
	defer wg.Done()
	close(ch)
}

func consumer(ch <-chan int, wg *sync.WaitGroup) {
	for num := range ch {
		fmt.Println("消費:", num)
	}
	defer wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int)
	go producer(ch, &wg)
	go consumer(ch, &wg)
	wg.Wait()
	fmt.Println("完成")
}
