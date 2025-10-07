package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var cache = &sync.Map{}

func calculateFibonacci(n int) *big.Int {
	if n <= 0 {
		return big.NewInt(0)
	}

	if n <= 1 {
		return big.NewInt(1)
	}

	a := big.NewInt(0)
	b := big.NewInt(1)
	var result *big.Int
	for i := 2; i <= n; i++ {
		result = new(big.Int).Set(a)
		result.Add(result, b)
		a.Set(b)
		b.Set(result)
	}

	return result
}

func fibonacciHandler(c *gin.Context) {
	n := c.DefaultQuery("n", "0")

	nInt, err := strconv.Atoi(n)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if cacheResult, ok := cache.Load(nInt); ok {
		fmt.Println("Cache Hit!")
		c.JSON(http.StatusOK, gin.H{"result": cacheResult})
		return
	}

	fmt.Println("Cache miss!")
	result := calculateFibonacci(nInt)
	cache.Store(nInt, result)

	go func() {
		time.Sleep(5 * time.Minute)
		cache.Delete(nInt)
	}()

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func main() {
	r := gin.Default()

	r.GET("/fibonacci", fibonacciHandler)

	port := ":8080"

	r.Run(port)

}
