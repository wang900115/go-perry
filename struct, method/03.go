package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) SayHello() string {
	return "你好，我是" + p.Name
}

func main() {
	p := Person{Name: "小明", Age: 25}
	fmt.Println(p.SayHello())
}
