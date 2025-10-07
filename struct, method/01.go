package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{Name: "小明", Age: 25}
	fmt.Println("姓名", p.Name, "年齡:", p.Age)
}
