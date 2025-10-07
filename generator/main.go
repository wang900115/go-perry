package main

import (
	"fmt"
	"math"
)

func nextPrime(n int) int {
	if n < 2 {
		return 2
	}

	for {
		n++
		if isPrime(n) {
			return n
		}
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func primeGenerator(done <-chan bool, operation func(int) int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		num := 1
		for {
			num = operation(num)
			select {
			case <-done:
				fmt.Println("Received done signal, stopping generator.")
				return
			case ch <- num:
			}
		}
	}()
	return ch
}

func main() {
	done := make(chan bool)

	i := 0
	for prime := range primeGenerator(done, nextPrime) {
		i++
		fmt.Println("Prime #", i, ": ", prime)
		if i == 10 {
			break
		}
	}

	done <- true
}
