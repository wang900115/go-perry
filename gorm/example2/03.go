package main

import (
	"fmt"
	"learn/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initializeDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=wang900115 dbname=learn port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil
	}
	return db
}

func main() {
	db := initializeDB()
	var actor model.Actor
	db.Where("name = ?", "Leonardo DiCaprio").Preload("Movies").First(&actor)
	fmt.Println(actor.ID, actor.Name)

	for _, movie := range actor.Movies {
		fmt.Println("Leonardo DiCaprio movies:", movie.Name)
	}
}
