// Package entity defines database table and its relations
package entity

import (
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
)

type ResponseAnalyzeCV struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:char(36)"`
	FileID   string    `json:"file_id" gorm:"type:varchar(256)"`
	Response string    `json:"response" gom:"type:nvarchar(4096)"`
}

func (a *ResponseAnalyzeCV) ParseToDTOResponseAnalyzeCV() dto.ResponseAnalyzeCV {
	return dto.ResponseAnalyzeCV{
		ID:       a.ID,
		FileID:   a.FileID,
		Response: a.Response,
	}
}
