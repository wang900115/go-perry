package model

type User2 struct {
	ID           uint    `gorm:"primaryKey"`
	Email        string  `gorm:"size:255"`
	UserSettings YamlMap `gorm:"type:text"`
}
