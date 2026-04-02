// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey"`
	OrderID     uuid.UUID      `gorm:"type:char(36);"`
	Token       string         `gorm:"type(36)"`
	RedirectURL string         `gorm:"type:nvarchar(256)"`
	CreatedAt   time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *Payment) ParseToDTOResponseCreateMidtransOrder() dto.ResponseCreateMidtransOrder {
	return dto.ResponseCreateMidtransOrder{
		Token:       p.Token,
		RedirectURL: p.RedirectURL,
	}
}
