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

func UserFromDomain(domain string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email like ?", "%"+domain)
	}
}

func CreditCardOrders(db *gorm.DB) *gorm.DB {
	return db.Where("payment_mode = ?", "CreditCard")
}

func main() {
	db := initializeDB()
	var users []model.User1

	db.Scopes(UserFromDomain("example.com")).Preload("Orders", CreditCardOrders).Find(&users)

	for _, user := range users {
		fmt.Println(user.ID, user.Name, user.Email)
		for _, order := range user.Orders {
			fmt.Printf(" Price: %d, Payment Mode: %s\n", order.Price, order.PaymentMode)
		}
	}

}
