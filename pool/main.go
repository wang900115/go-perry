package main

/// Greate for short-live objects
/// Objects in the poll may be garbage collected

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	allocationCnt := 0
	pool := sync.Pool{
		New: func() interface{} {
			fmt.Print(".")
			allocationCnt++
			return make([]byte, 1024)
		},
	}

	// obj := pool.Get().([]byte)
	// fmt.Printf("Got object from pool of length: %d\n", len(obj))

	// pool.Put(obj)

	// reusedObj := pool.Get().([]byte)
	// fmt.Printf("Got reused object from pool of length: %d\n", len(reusedObj))

	// pool.Put(reusedObj)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			obj := pool.Get().([]byte)
			fmt.Print("-")
			time.Sleep(100 * time.Millisecond)
			pool.Put(obj)
			wg.Done()
		}()
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()

	fmt.Println("\n Number of allocations: ", allocationCnt)
}
