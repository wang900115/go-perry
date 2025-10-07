package model

import (
	"time"

	"gorm.io/gorm"
)

type User1 struct {
	gorm.Model
	Name   string
	Email  string
	Orders []Order
}

type Order struct {
	gorm.Model
	User1ID     uint
	OrderTime   time.Time
	PaymentMode string
	Price       int
	User1       User1
}
