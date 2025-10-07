package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, result chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d 開始處理任務 %d \n", id, j)
		time.Sleep(10 * time.Second)
		fmt.Printf("Worker %d 完成任務 %d\n", id, j)
		result <- j * 2
	}
}

func main() {
	const numJobs = 5
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println("收到結果:", r)
	}
	fmt.Println("全部完成")
}
