package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	person := Person{Name: "Alan", Age: 30}
	val := reflect.ValueOf(person)

	nameField := val.FieldByName("Name")
	fmt.Printf("Name: %v\n", nameField) // Name: "Alan"

	ageField := val.FieldByName("Age")
	fmt.Printf("Age: %v\n", ageField) // Age: 30

	valPtr := reflect.ValueOf(&person).Elem()
	valPtr.FieldByName("Age").SetInt(15)
	fmt.Printf("Updated Age: %v\n", person.Age) // Updated Age: 15
}
