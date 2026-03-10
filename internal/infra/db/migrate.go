package db

import (
	"log"

	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		entity.User{},
		entity.UserDetail{},
		entity.Course{},
		entity.Order{},
		entity.Payment{},
	)
	if err != nil {
		log.Panic("database migration failed")
	}

	log.Println("database migration complete")
}
