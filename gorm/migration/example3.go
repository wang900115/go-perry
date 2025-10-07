package main

import (
	"fmt"
	"learn/model"
	"time"

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
	db.AutoMigrate(&model.User1{}, &model.Order{})

	users := []model.User1{
		{
			Name:  "Alice",
			Email: "alice@example.com",
			Orders: []model.Order{
				{OrderTime: time.Now(), PaymentMode: "CreditCard", Price: 1200},
				{OrderTime: time.Now().Add(-48 * time.Hour), PaymentMode: "Cash", Price: 800},
			},
		},
		{
			Name:  "Bob",
			Email: "bob@example.com",
			Orders: []model.Order{
				{OrderTime: time.Now().Add(-24 * time.Hour), PaymentMode: "CreditCard", Price: 1500},
			},
		},
		{
			Name:  "Charlie",
			Email: "charlie@example.com",
			Orders: []model.Order{
				{OrderTime: time.Now(), PaymentMode: "CreditCard", Price: 999},
				{OrderTime: time.Now().Add(-72 * time.Hour), PaymentMode: "Cash", Price: 500},
				{OrderTime: time.Now().Add(-24 * time.Hour), PaymentMode: "CreditCard", Price: 600},
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
