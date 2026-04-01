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
		entity.UserContact{},
		entity.UserEducation{},
		entity.UserLanguage{},
		entity.UserEmployment{},
		entity.UserSeniority{},
		entity.UserWorkExperience{},
		entity.UserHardSkill{},
		entity.UserSoftSkill{},
		entity.UserTools{},
		entity.UserLink{},
		entity.UserSubscription{},
		entity.Course{},
		entity.ResponseAnalyzeCV{},
		entity.Order{},
		entity.Payment{},
	)
	if err != nil {
		log.Panic("database migration failed")
	}

	log.Println("database migration complete")
}
