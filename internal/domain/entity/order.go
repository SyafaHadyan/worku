// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey"`
	UserID       uuid.UUID      `gorm:"type:char(36);"`
	DurationDays int            `gorm:"type:integer"`
	Status       string         `gorm:"type:nvarchar(128)"`
	CreatedAt    time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;autoUpdateTime"`
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
