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
	db.AutoMigrate(&model.Movie{}, &model.Actor{})

	movies := []model.Movie{
		{Name: "Inception"},
		{Name: "Titanic"},
		{Name: "The Matrix"},
		{Name: "The Dark Knight"},
		{Name: "Interstellar"},
	}

	for _, movie := range movies {
		if err := db.Create(&movie).Error; err != nil {
			fmt.Println("Error inserting movie:", err)
		} else {
			fmt.Printf("Inserted movie: %s\n", movie.Name)
		}
	}

	actors := []model.Actor{
		{Name: "Leonardo DiCaprio"},
		{Name: "Keanu Reeves"},
		{Name: "Christian Bale"},
		{Name: "Joseph Gordon-Levitt"},
		{Name: "Matthew McConaughey"},
	}

	for _, actor := range actors {
		if err := db.Create(&actor).Error; err != nil {
			fmt.Println("Error inserting actor:", err)
		} else {
			fmt.Printf("Inserted actor: %s\n", actor.Name)
		}
	}
}
