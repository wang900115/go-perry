package main

import (
	"fmt"
	"strconv"
)

// string -> number
func main() {
	var (
		err error
		s1  string = "5"
		s2  string = "abc"
		i1  int
	)

	i1, err = strconv.Atoi(s1)
	fmt.Printf("%T %v \n", i1, i1)

	_, err = strconv.Atoi(s2)
	if err != nil {
		fmt.Print(err)
	}
}
