package database

import "github.com/jinzhu/gorm"

var (
	UserTable  = new(User)
	ImageTable = new(Image)
)

type User struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null"`
	Pass    string `gorm:"type:varchar(64);not null"`
	Salt    string `gorm:"type:varchar(64);not null"`
	IconURL string
	Images  []Image
}

type Image struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null"`
	Ext         string `gorm:"type:varchar(8);not null"`
	Hash        string `gorm:"type:varchar(255);not null"`
	OriginalURL string `gorm:"type:varchar(255);not null"`
}

func IsExistUser(db *gorm.DB, username string) bool {
	return !db.Where("name = ?", username).First(UserTable).RecordNotFound()
}
