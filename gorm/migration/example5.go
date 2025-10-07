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
	db.AutoMigrate(&model.User2{})

	user := model.User2{
		Email: "abc@example.com",
		UserSettings: model.YamlMap{
			"theme":         "dark",
			"notifications": true,
		},
	}

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("fail create user2", err)
	} else {
		fmt.Println("inserted user2: ", user.ID)
	}
}
