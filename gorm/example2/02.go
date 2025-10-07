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
	var movie model.Movie
	db.Where("name = ?", "Titanic").Preload("Actors").First(&movie)
	fmt.Println(movie.ID, movie.Name)

	for _, actor := range movie.Actors {
		fmt.Println("Inception actor:", actor.Name)
	}
}
