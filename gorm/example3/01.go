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

func PriceGreaterThan30(db *gorm.DB) *gorm.DB {
	return db.Where("price > ?", 30)
}

func CreditCardOrders(db *gorm.DB) *gorm.DB {
	return db.Where("payment_mode = ?", "CreditCard")
}

func main() {
	db := initializeDB()
	var orders []model.Order

	db.Scopes(CreditCardOrders, PriceGreaterThan30).Find(&orders)
	for _, order := range orders {
		fmt.Printf("Price: %d, Payment Mode: %s\n", order.Price, order.PaymentMode)
	}

}
