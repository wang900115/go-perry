package main

import (
	"fmt"
	"sync"
)

func managerTickets(ticketChannel chan int, doneChannel chan bool, ticket *int) {
	for {
		select {
		case user := <-ticketChannel:
			if *ticket > 0 {
				*ticket--
				fmt.Printf(" User %d purchased the ticket, remaining %d tickets \n", user, *ticket)
			} else {
				fmt.Printf(" User %d found no ticket \n", user)
			}
		case <-doneChannel:
			fmt.Printf("Ticket remaing %d", *ticket)
		}
	}

}

func buyTicket(wg *sync.WaitGroup, ticketChannel chan int, userID int) {
	defer wg.Done()
	ticketChannel <- userID
}

func main() {

	var wg sync.WaitGroup
	tickets := 200
	users := 2000
	ticketChannel := make(chan int)
	doneChannel := make(chan bool)

	go managerTickets(ticketChannel, doneChannel, &tickets)

	for userID := 0; userID < users; userID++ {
		wg.Add(1)
		go buyTicket(&wg, ticketChannel, userID)
	}

	wg.Wait()
	doneChannel <- true
}
