package main

import (
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	url := "http://localhost:8080/"
	for i := 1; i <= 50; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making request", err)
			continue
		}

		fmt.Printf("Request: %2d: Status: %d\n", i, resp.StatusCode)
		time.Sleep(100 * time.Millisecond)
	}
}
