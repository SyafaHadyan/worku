// Package entity defines database table and its relations
package entity

import (
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
)

type ResponseAnalyzeCV struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID   uuid.UUID `gorm:"type:char(36)"`
	FileID   string    `gorm:"type:varchar(256)"`
	Response string    `gorm:"type:mediumtext"`
}

func (a *ResponseAnalyzeCV) ParseToDTOResponseAnalyzeCV() dto.ResponseAnalyzeCV {
	return dto.ResponseAnalyzeCV{
		ID:       a.ID,
		FileID:   a.FileID,
		Response: a.Response,
	}
}
