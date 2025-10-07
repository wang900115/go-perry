package main

// number -> string
import (
	"fmt"
	"strconv"
)

func main() {
	var str1 string
	var i int = 10
	fmt.Printf("%T %v", i, i)

	// strconv.Itoa
	str1 = strconv.Itoa(i)
	fmt.Printf("%T %v", str1, str1)

	// fmt.Sprintf
	str2 := fmt.Sprintf("%d", i)
	fmt.Printf("%T %v", str2, str2)
}
