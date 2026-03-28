// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID     uuid.UUID      `json:"order_id" gorm:"type:char(36);"`
	Token       string         `json:"token" gorm:"type(36)"`
	RedirectURL string         `json:"redirect_url" gorm:"type:nvarchar(256)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *Payment) ParseToDTOResponseCreateMidtransOrder() dto.ResponseCreateMidtransOrder {
	return dto.ResponseCreateMidtransOrder{
		Token:       p.Token,
		RedirectURL: p.RedirectURL,
	}
}
