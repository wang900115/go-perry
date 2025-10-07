package main

import "fmt"

type Coffee interface {
	Cost() float64
	Description() string
}

type SimpleCoffee struct{}

func (s *SimpleCoffee) Cost() float64 {
	return 2.0
}

func (s *SimpleCoffee) Description() string {
	return "Simple Coffee"
}

// Decorator 1: Milk

type Milk struct {
	coffee Coffee
}

func (m *Milk) Cost() float64 {
	return m.coffee.Cost() + 0.5
}

func (m *Milk) Description() string {
	return m.coffee.Description() + ", Milk"
}

// Decorator 2: Caramel

type Caramel struct {
	coffee Coffee
}

func (c *Caramel) Cost() float64 {
	return c.coffee.Cost() + 1
}

func (c *Caramel) Description() string {
	return c.coffee.Description() + ", Caramel"
}

func main() {
	coffee := &SimpleCoffee{}

	coffeeWithMilk := &Milk{coffee: coffee}
	fmt.Println("Coffee with Milk: ", coffeeWithMilk.Description(), "Cost: ", coffeeWithMilk.Cost())

	coffeeWithCaramel := &Caramel{coffee: coffee}
	fmt.Println("Coffee with Caramel: ", coffeeWithCaramel.Description(), "Cost: ", coffeeWithCaramel.Cost())

	coffeeWithMilkAndCaramel := &Milk{coffee: &Caramel{coffee: coffee}}
	fmt.Println("Coffee with Milk and Caramel:", coffeeWithMilkAndCaramel.Description(), "Cost: ", coffeeWithMilkAndCaramel.Cost())
}
