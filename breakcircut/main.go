package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type State int

const (
	Open State = iota
	Closed
	HalfOpen
)

type Breaker struct {
	mu    sync.Mutex
	state State
	// records
	failureCount int
	successCount int
	// configuration
	failThreshold    int
	successThreshold int
	resetTimeout     time.Duration
	// timing records
	lastFailTime time.Time
}

func NewBreaker(failThreshold, successThreshold int, resetTimeout time.Duration) *Breaker {
	return &Breaker{
		state:            Open,
		failThreshold:    failThreshold,
		successThreshold: successThreshold,
		resetTimeout:     resetTimeout,
	}
}

func (b *Breaker) goToOpen() {
	b.state = Open
	b.failureCount = 0
	b.successCount = 0
	b.lastFailTime = time.Now()
}

func (b *Breaker) goToClosed() {
	b.state = Closed
	b.failureCount = 0
	b.successCount = 0
}

func (b *Breaker) goToHalfOpen() {
	b.state = HalfOpen
	b.failureCount = 0
	b.successCount = 0
}

func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Open:
		if time.Since(b.lastFailTime) > b.resetTimeout {
			b.goToHalfOpen()
			return nil
		}
		return fmt.Errorf("circuit breaker is open, rejecting all requests")
	default:
		return nil
	}
}

func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case HalfOpen:
		b.successCount++
		if b.successCount >= b.successThreshold {
			b.goToClosed()
		}
	case Closed:
		b.failureCount = 0
	}
}

func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case HalfOpen:
		b.goToOpen()
	case Closed:
		b.failureCount++
		if b.failureCount >= b.failThreshold {
			b.goToOpen()
		}
	}
}

func unreliable(pFail float64) error {
	if rand.Float64() < pFail {
		return fmt.Errorf("unreliable service failed")
	}
	time.Sleep(time.Duration(rand.IntN(50)) * time.Millisecond)
	return nil
}
