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
	var notes []model.Note
	db := initializeDB()
	db.Where("user_id = ?", 1).Order("title DESC").Find(&notes)
	for _, note := range notes {
		fmt.Println(note.ID, note.Title, note.Content, note.UserID)
	}
}
