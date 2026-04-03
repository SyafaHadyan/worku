// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"gorm.io/gorm"
)

type AIDBItf interface {
	ResponseAnalyzeCV(responseAnalyzeCV *entity.ResponseAnalyzeCV) error
}

type AIDB struct {
	db *gorm.DB
}

func NewAIDB(db *gorm.DB) AIDBItf {
	return &AIDB{
		db: db,
	}
}

func (r *AIDB) ResponseAnalyzeCV(responseAnalyzeCV *entity.ResponseAnalyzeCV) error {
	return r.db.
		Create(responseAnalyzeCV).
		Error
}
