// Package main demonstrates time.Timer, time.Ticker and time.AfterFunc,
// with a focus on the Go 1.23+ semantics of Stop and Reset.
package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	part1TimeoutInLoop()
	part2StopResetNoStaleValue()
	part3TickerDropsTicks()
	part4AfterFunc()
	part5LenAndCap()
}

// Part 1: a timeout inside a loop. time.After creates a new timer on every
// iteration; one reusable Timer with Reset does the same job.
func part1TimeoutInLoop() {
	fmt.Println("== Part 1: timeout in a loop ==")
	events := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(20 * time.Millisecond)
			events <- i
		}
		// then go quiet so the idle timeout fires
	}()

	const idle = 100 * time.Millisecond
	t := time.NewTimer(idle)
	defer t.Stop()
	for {
		select {
		case e := <-events:
			fmt.Println("event", e)
			t.Reset(idle) // Go 1.23+: safe without draining t.C first
		case <-t.C:
			fmt.Println("idle timeout, leaving loop")
			return
		}
	}
}

// Part 2: since Go 1.23, after Stop or Reset returns, no stale value
// from the old timer can be received from t.C.
func part2StopResetNoStaleValue() {
	fmt.Println("== Part 2: Stop/Reset leave no stale value ==")

	t := time.NewTimer(10 * time.Millisecond)
	time.Sleep(30 * time.Millisecond) // timer has expired, nobody received

	t.Reset(200 * time.Millisecond)
	start := time.Now()
	<-t.C
	waited := time.Since(start)
	// Before Go 1.23 the stale value would be received immediately (~0ms).
	fmt.Println("after Reset, waited >= 150ms:", waited >= 150*time.Millisecond)

	t2 := time.NewTimer(10 * time.Millisecond)
	time.Sleep(30 * time.Millisecond)
	stopped := t2.Stop()
	fmt.Println("Stop on expired-but-unreceived timer returns:", stopped)
	select {
	case <-t2.C:
		fmt.Println("received stale value (old behaviour)")
	default:
		fmt.Println("t2.C is empty after Stop")
	}
}

// Part 3: a Ticker never queues more than one tick. A slow receiver
// loses ticks instead of building up a backlog.
func part3TickerDropsTicks() {
	fmt.Println("== Part 3: Ticker drops ticks for a slow receiver ==")
	const period = 20 * time.Millisecond
	tk := time.NewTicker(period)
	defer tk.Stop()

	deadline := time.Now().Add(300 * time.Millisecond)
	received := 0
	for time.Now().Before(deadline) {
		<-tk.C
		received++
		time.Sleep(70 * time.Millisecond) // slower than the ticker
	}
	// 300ms / 20ms = 15 ticks were produced, but only about 4-5 are
	// received: missed ticks are dropped, not queued.
	fmt.Println("ticks produced ~15, received < 8:", received < 8)
}

// Part 4: AfterFunc runs f in its own goroutine. t.C is nil.
// Stop returns false if f has already started.
func part4AfterFunc() {
	fmt.Println("== Part 4: AfterFunc ==")
	var fired atomic.Int32

	a := time.AfterFunc(50*time.Millisecond, func() { fired.Add(1) })
	fmt.Println("AfterFunc timer channel is nil:", a.C == nil)
	fmt.Println("Stop before firing returns:", a.Stop())

	b := time.AfterFunc(10*time.Millisecond, func() { fired.Add(10) })
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Stop after firing returns:", b.Stop())
	fmt.Println("fired total:", fired.Load()) // only b ran
}

// Part 5: since Go 1.23 the timer channel is synchronous (unbuffered):
// len and cap report 0, even after the timer has expired.
func part5LenAndCap() {
	fmt.Println("== Part 5: timer channel is unbuffered ==")
	t := time.NewTimer(5 * time.Millisecond)
	defer t.Stop()
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("len=%d cap=%d\n", len(t.C), cap(t.C))
	<-t.C
	fmt.Println("but a receive still gets the value")
}
