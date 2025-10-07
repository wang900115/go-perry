package main

import "fmt"

type foo struct {
	i int
}

type bar foo

func main() {
	var b bar = bar{1}
	var f foo = foo(b)

	fmt.Printf("%T %v", b, b)
	fmt.Printf("%T %v", f, f)

}
