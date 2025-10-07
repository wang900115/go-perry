package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println("前三個", s[:3])
	fmt.Println("後三個", s[2:])
}
