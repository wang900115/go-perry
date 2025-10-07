package main

import "fmt"

func add(a int, b int) int {
	return a + b
}

func main() {
	result := add(3, 6)
	fmt.Println("3 + 6 =", result)
}
