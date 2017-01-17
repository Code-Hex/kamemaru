package database

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Pass    string `gorm:"not null"`
	IconURL string
	Images  []Image
}

type Image struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Ext         string `gorm:"type:varchar(8);not null"`
	Hash        string `gorm:"not null"`
	OriginalURL string `gorm:"not null"`
}
