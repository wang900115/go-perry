package main

import (
	"fmt"
	"sort"
)

type User struct {
	Name string
	Age  int
}

func main() {
	users := []User{
		{
			Name: "Rage",
			Age:  28,
		},
		{
			Name: "Mino",
			Age:  31,
		},
		{
			Name: "Joda",
			Age:  12,
		},
		{
			Name: "Kap",
			Age:  58,
		},
		{
			Name: "loam",
			Age:  44,
		},
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Age < users[j].Age
	})

	fmt.Println("Sorted users by age: ", users)

	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Age > users[j].Age
	})

	fmt.Println("Stabled users reversed by age: ", users)
}
