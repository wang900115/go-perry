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
	db.AutoMigrate(&model.User{}, &model.Note{}, &model.CreditCard{})

	users := []model.User{
		{
			Name: "John Doe",
			Notes: []model.Note{
				{Title: "Note A-1", Content: "Content A-1"},
				{Title: "Note A-2", Content: "Content A-2"},
			},
			CreditCard: model.CreditCard{
				Number: "1234-5678-9012-3456",
			},
		},
		{
			Name: "Jane Smith",
			Notes: []model.Note{
				{Title: "Note B-1", Content: "Content B-1"},
				{Title: "Note B-2", Content: "Content B-2"},
			},
			CreditCard: model.CreditCard{
				Number: "9876-5432-1098-7654",
			},
		},
		{
			Name: "Alice Johnson",
			Notes: []model.Note{
				{Title: "Note C-1", Content: "Content C-1"},
				{Title: "Note C-2", Content: "Content C-2"},
			},
			CreditCard: model.CreditCard{
				Number: "5555-4444-3333-2222",
			},
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			fmt.Println("Error inserting user:", err)
		} else {
			fmt.Printf("Inserted user: %s\n", user.Name)
		}
	}
}
