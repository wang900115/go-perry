package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 有一個input channel 一旦 done signal 出現就停止接受訊息
func orDone(done <-chan struct{}, input <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-input:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

// 從 input channel 中讀取資料 複製到兩個channel
func tee(input <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)
		for v := range input {
			var o1, o2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case o1 <- v:
					o1 = nil
				case o2 <- v:
					o2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func transactionGenerator(done <-chan struct{}) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 0; i < 200; i++ {
			select {
			case <-done:
				fmt.Println("Transaction Generator: Shutting down ...")
				return
			case out <- rand.Intn(1000):
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	return out
}

func fraudDetectionPipeline(done <-chan struct{}, input <-chan int) {
	for transaction := range orDone(done, input) {
		if transaction > 700 {
			fmt.Println("Fraud Detection: Suspicious transaction detected:", transaction)
		} else {
			fmt.Println("Fraud Detection: Normal transaction:", transaction)
		}
	}

	fmt.Println("Fraud Detection Pipeline shut down")
}

func analyticsPipeline(done <-chan struct{}, input <-chan int) {
	total := 0
	count := 0
	for transaction := range orDone(done, input) {
		fmt.Println("Analytics: Processing transaction:", transaction)
		total += transaction
		count++
	}

	if count > 0 {
		fmt.Printf("Analytics: Average transaction amount: %d\n", total/count)
	}
	fmt.Println("Analytic Pipeline shut down")
}

func main() {
	done := make(chan struct{})

	source := transactionGenerator(done)
	fraudCh, analyticsCh := tee(source)

	go fraudDetectionPipeline(done, fraudCh)
	go analyticsPipeline(done, analyticsCh)

	time.Sleep(10 * time.Second)
	close(done)

	time.Sleep(1 * time.Second)
	fmt.Println("Main: All pipeline stopped")
}
