package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Job struct {
	ID int
}

func worker(id int, jobs <-chan Job, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d 處理任務 %d\n", id, job.ID)
		time.Sleep(1 * time.Second) // 模擬耗時
		results <- fmt.Sprintf("任務 %d 完成", job.ID)
	}
}

func main() {
	r := gin.Default()
	jobs := make(chan Job, 10)
	results := make(chan string, 10)
	var wg sync.WaitGroup

	// 啟動 3 個 worker
	const numWorkers = 3
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, &wg)
	}

	// HTTP 端點：提交任務
	r.GET("/process/:id", func(c *gin.Context) {
		id := c.Param("id")
		jobID := 0
		fmt.Sscanf(id, "%d", &jobID)
		jobs <- Job{ID: jobID}
		c.JSON(200, gin.H{"message": fmt.Sprintf("任務 %d 已提交", jobID)})
	})

	// 收集結果（模擬後台處理）
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	// 關閉 channel 和等待 worker
	go func() {
		wg.Wait()
		close(results)
	}()

	// 啟動服務
	r.Run(":8080")
}
