package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p *Person) GrowUp() {
	p.Age++
}

func main() {
	p := Person{Name: "小明", Age: 25}
	p.GrowUp()
	fmt.Println(p.Name, "現在", p.Age, "歲")
}
