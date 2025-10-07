package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex

func buyticket(wg *sync.WaitGroup, userID int, tickets *int) {
	defer wg.Done()
	mutex.Lock()
	if *tickets > 0 {
		*tickets--
		fmt.Printf("User %d Buy ticket, remaining %d \n", userID, *tickets)
	} else {
		fmt.Printf("User %d No ticket can buy \n", userID)
	}
	mutex.Unlock()
}

func main() {
	var tickets int = 500
	var wg sync.WaitGroup
	for userID := 0; userID < 2000; userID++ {
		wg.Add(1)
		go buyticket(&wg, userID, &tickets)
	}
	wg.Wait()
}
