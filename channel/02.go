package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("goroutine 執行中")
	}()

	wg.Wait()
	fmt.Println("主執行緒結束")
}
