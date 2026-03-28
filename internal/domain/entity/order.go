// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:char(36);"`
	DurationDays int            `json:"duration" gorm:"type:integer"`
	Status       string         `json:"status" gorm:"type:nvarchar(128)"`
	CreatedAt    time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (o *Order) ParseToDTOResponseCreateOrder() dto.ResponseCreateOrder {
	return dto.ResponseCreateOrder{
		ID:           o.ID,
		UserID:       o.UserID,
		DurationDays: o.DurationDays,
		Status:       o.Status,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

func (o *Order) ParseToDTOResponseGetOrderInfo() dto.ResponseGetOrderInfo {
	return dto.ResponseGetOrderInfo{
		ID:           o.ID,
		UserID:       o.UserID,
		DurationDays: o.DurationDays,
		Status:       o.Status,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

func (o *Order) ParseToDTOResponseGetOrderList() dto.ResponseGetOrderList {
	return dto.ResponseGetOrderList{
		ID:           o.ID,
		UserID:       o.UserID,
		DurationDays: o.DurationDays,
		Status:       o.Status,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}
