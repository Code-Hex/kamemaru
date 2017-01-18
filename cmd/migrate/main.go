package main

import (
	"github.com/Code-Hex/kamemaru/internal/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "host=localhost dbname=kamemaru sslmode=disable")
	if err != nil {
		panic(err)
	}

	user, image := database.UserTable, database.ImageTable
	db.DropTableIfExists(user, image)
	db.AutoMigrate(user, image)
	defer db.Close()
}
