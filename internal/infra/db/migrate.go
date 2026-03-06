package db

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		log.Panic("database migration failed")
	}

	log.Println("database migration complete")
}
