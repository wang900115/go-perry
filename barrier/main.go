package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var data int
var ready atomic.Bool
var mu sync.Mutex

// 1: using atomic store to ensure that the write to data is visible to other goroutines when ready is set to true
// release barrier: when ready is set to true, all previous writes (data = 42) are visible to other goroutines

// 2. using atomic load to ensure that the read of ready is synchronized with the write to ready in the other goroutine
// acquire barrier: when ready is read as true, all previous writes (data = 42) are visible to this goroutine
func main() {
	go func() {
		data = 42
		ready.Store(true)
	}()

	for !ready.Load() {
		// spin
	}
	fmt.Println(data)

}

// 1. using mutex to ensure that the write to data is visible to other goroutines when the lock is released
// release barrier: when the lock is released, all previous writes (data = 42) are visible to other goroutines

// 2. using mutex to ensure that the read of data is synchronized with the write to data in the other goroutine
// acquire barrier: when the lock is acquired, all previous writes (data = 42) are visible to this goroutine
func main2() {
	go func() {
		mu.Lock()
		data = 42
		mu.Unlock()
	}()

	mu.Lock()
	fmt.Println(data)
	mu.Unlock()
}

// 1. using channel to ensure that the write to data is visible to other goroutines when the value is sent on the channel
// release barrier: when the value is sent on the channel, all previous writes (data = 42) are visible to other goroutines

// 2. using channel to ensure that the read of data is synchronized with the write to data in the other goroutine
// acquire barrier: when the value is received from the channel, all previous writes (data = 42) are visible to this goroutine
func main3() {

	ch := make(chan int)
	go func() {
		data := 42
		ch <- data
	}()
	v := <-ch
	fmt.Println(v)
}

// func main() {
// 	go func() {
// 		data = 42
// 		ready = true
// 	}()

// 	for !ready {
// 		// spin
// 	}

// 	fmt.Println(data)
// 	time.Sleep(1 * time.Second)
// }
