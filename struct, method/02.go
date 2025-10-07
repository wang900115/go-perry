package main

import "fmt"

type Address struct {
	City    string
	Country string
}

type Person struct {
	Name    string
	Age     int
	Address Address
}

func main() {
	p := Person{
		Name: "小明",
		Age:  25,
		Address: Address{
			City:    "台北",
			Country: "台灣",
		},
	}
	fmt.Println(p.Name, "住在", p.Address.City)
}
