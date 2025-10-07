package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Tags        []string               `gorm:"type:bytes;serializer:json"`
	Spec        map[string]interface{} `gorm:"serializer:json"`
	SpecGob     map[string]interface{} `gorm:"type:bytes;serializer:gob"`
	CreatedTime int64                  `gorm:"serializer:unixtime;type:time"`
}
