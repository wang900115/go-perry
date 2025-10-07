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

	db.AutoMigrate(&model.Product{})
	createdAt := time.Now()

	data := model.Product{
		Name: "Apple iPhone 13",
		Tags: []string{"smartphone", "iphone", "apple", "cell phone", "5g", "camera", "retina display"},
		Spec: map[string]interface{}{
			"name":       "Apple iPhone 13",
			"display":    "6.1 inches",
			"resolution": "2532 x 1170 pixels",
			"processor":  "Apple A15 Bionic",
			"ram":        "6GB",
			"storage":    "128GB",
		},
		SpecGob: map[string]interface{}{
			"name":       "Apple iPhone 13",
			"display":    "6.1 inches",
			"resolution": "2532 x 1170 pixels",
			"processor":  "Apple A15 Bionic",
			"ram":        "6GB",
			"storage":    "128GB",
		},
		CreatedTime: createdAt.Unix(),
	}

	if err := db.Create(&data).Error; err != nil {
		fmt.Println("Error inserting product:", err)
	} else {
		fmt.Println("Inserted product: \n", data.Name)
	}

}
