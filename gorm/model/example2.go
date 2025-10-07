package model

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Name   string
	Actors []Actor `gorm:"many2many:filmographs;"`
}

type Actor struct {
	gorm.Model
	Name   string
	Movies []Movie `gorm:"many2many:filmographs;"`
}
