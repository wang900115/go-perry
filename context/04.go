package main

import (
	"context"
	"fmt"
)

type contextKey string

const userIDKey = contextKey("userID")

func main() {
	ctx := context.WithValue(context.Background(), userIDKey, 42)
	process(ctx)
}

func process(ctx context.Context) {
	userId := ctx.Value(userIDKey)
	fmt.Println("User ID from context:", userId)
}
