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
	var users []model.User

	db := initializeDB()

	db.Preload("CreditCard").Preload("Notes").Find(&users)

	for _, user := range users {
		fmt.Println(user.ID, user.Name, user.Notes, user.CreditCard.Number)
	}

}
