package main

import (
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer runtime.UnlockOSThread()
		runtime.LockOSThread()
		for i := 0; i < 5; i++ {
			println("[1] This goroutine is locked to OS thread:", getOSThreadID())
			time.Sleep(time.Second * 2)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			println("[2] This goroutine is NOT locked to OS thread:", getOSThreadID())
			time.Sleep(time.Second * 2)
		}
	}()
	wg.Wait()
}

func getOSThreadID() int {
	return syscall.Gettid()
}
