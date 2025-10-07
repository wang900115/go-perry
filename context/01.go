package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	fmt.Println("This is the root context:", ctx)
}
