package main

import "fmt"

type Oven interface {
	Heat() string
}
type Ingredient interface {
	Mix() string
}
type GasOven struct{}

func (gasOven GasOven) Heat() string {
	return "Heating with Gas Oven!"
}

type ElectricOven struct{}

func (electricOven ElectricOven) Heat() string {
	return "Heating with electrice oven!"
}

type Flour struct{}

func (f Flour) Mix() string {
	return "Mixing flour!"
}

type Sugar struct{}

func (s Sugar) Mix() string {
	return "Mixing sugar!"
}

type Butter struct{}

func (b Butter) Mix() string {
	return "Mixing butter!"
}

type Bakery struct {
	oven        Oven
	ingredients []Ingredient
}

func (b *Bakery) Bake() {
	fmt.Println(b.oven.Heat())
	for _, ingredient := range b.ingredients {
		fmt.Println(ingredient.Mix())
	}

	fmt.Println("Baking an asesome pastry!")
}

func main() {

	gasOven := GasOven{}
	electricOven := ElectricOven{}

	flour := Flour{}
	sugar := Sugar{}
	butter := Butter{}

	bakery := Bakery{oven: gasOven, ingredients: []Ingredient{flour, sugar}}
	bakery.Bake()
	bakery = Bakery{oven: electricOven, ingredients: []Ingredient{flour, butter}}
	bakery.Bake()

}
