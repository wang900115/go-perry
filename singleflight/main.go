package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var count int = 0
var lock sync.Mutex = sync.Mutex{}

func getData() (interface{}, error) {
	lock.Lock()
	count++
	lock.Unlock()
	return http.Get("https://google.com")
}

func getDataSingleFlight(g *singleflight.Group, wg *sync.WaitGroup) error {
	defer wg.Done()
	v, err, shared := g.Do("get-google", getData)
	if err != nil {
		return err
	}

	res := v.(*http.Response)
	fmt.Printf("status: %d, shared: %v\n", res.StatusCode, shared)
	return nil
}

func main() {
	var wg sync.WaitGroup
	var g singleflight.Group
	wg.Add(10)
	for i := 0; i < 10; i++ {
		if i == 4 {
			time.Sleep(1 * time.Second)
		}
		go getDataSingleFlight(&g, &wg)
	}
	wg.Wait()
	fmt.Printf("original function was called %d times", count)
}
