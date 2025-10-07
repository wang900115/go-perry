package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

var urls = []string{
	"https://www.codeheim.io",
	"https://golang.org",
	"https://pkg.go.dev/golang.org/x/sync/errgroup",
	"https://invalid-url",
}

func fetchPage(url string, mu *sync.Mutex, response *map[string]string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch %s: %s\n", url, err)
		return fmt.Errorf("failed to fetch %s: %w", url, err)
	}

	defer resp.Body.Close()

	fmt.Printf("Successfully fetched %s\n", url)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response from %s: %w", url, err)
	}

	mu.Lock()
	(*response)[url] = string(body)
	mu.Unlock()

	fmt.Printf("Successfully fetched response body of %s\n", url)

	return nil
}

func main() {
	var g errgroup.Group

	g.SetLimit(2)

	response := make(map[string]string)
	var mu sync.Mutex

	for _, url := range urls {
		g.Go(func() error {
			return fetchPage(url, &mu, &response)
		})

	}

	if err := g.Wait(); err != nil {
		fmt.Println("Error occured: ", err)
	} else {
		fmt.Println("All URLS fetched successfully")

		for url, content := range response {
			fmt.Printf("Response from %s: %s\n", url, content[:100])
		}
	}

}
