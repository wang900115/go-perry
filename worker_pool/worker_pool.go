package main

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Process()
}

type EmailTask struct {
	Email       string
	Subject     string
	MessageBody string
}

func (t *EmailTask) Process() {
	fmt.Println("Sending Email to : ", t.Email)
	time.Sleep(2 * time.Second)
}

type ImageProcessingTask struct {
	ImageUrl string
}

func (t *ImageProcessingTask) Process() {
	fmt.Println("Processing the image : ", t.ImageUrl)
	time.Sleep(5 * time.Second)
}

type WorkerPool struct {
	Tasks       []Task
	concurrency int
	tasksChan   chan Task
	wg          sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for task := range wp.tasksChan {
		task.Process()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run() {
	wp.tasksChan = make(chan Task, len(wp.Tasks))

	for i := 0; i < wp.concurrency; i++ {
		go wp.worker()
	}

	wp.wg.Add(len(wp.Tasks))

	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}

	close(wp.tasksChan)

	wp.wg.Wait()
}
