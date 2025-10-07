package main

import "fmt"

func main() {
	// var a byte = 10
	// var n int = int(a)

	// fmt.Println(n)

	a := 0x1234
	b := 1234.56
	c := 256

	fmt.Printf("%x \n", uint8(a))
	fmt.Printf("%d \n", int(b))
	fmt.Printf("%f \n", float64(c))

}
