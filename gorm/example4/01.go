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

	var result model.Product
	db.First(&result, "id = ?", 1)

	fmt.Println("Name: ", result.Name)
	fmt.Println("Tag: ", result.Tags)
	fmt.Println("Spec: ", result.Spec)
	fmt.Println("SpecGob: ", result.SpecGob)
	fmt.Println("CreatedTime: ", result.CreatedTime)
}
