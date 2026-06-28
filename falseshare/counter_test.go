package main_test

import (
	main "falsesharee"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkFalseSharing(b *testing.B) {
	var counter main.Counter

	b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			atomic.AddInt64(&counter.A, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			atomic.AddInt64(&counter.B, 1)
		}
	}()

	wg.Wait()

	runtime.KeepAlive(counter)
}
func BenchmarkPaddedCounter(b *testing.B) {
	counter := main.CustomCounter{}
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		for i := 0; i < b.N; i++ {
			atomic.AddInt64(&counter.A, 1)
		}
	}()
	go func() {
		defer wg.Done()

		for i := 0; i < b.N; i++ {
			atomic.AddInt64(&counter.B, 1)
		}
	}()
	wg.Wait()
	runtime.KeepAlive(counter)
}
