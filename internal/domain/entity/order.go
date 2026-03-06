// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	CanteenID uuid.UUID      `json:"canteen_id" gorm:"type:char(36);"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36);"`
	Quantity  uint32         `json:"quantity" gorm:"type:integer unsigned"`
	Status    string         `json:"status" gorm:"type:varchar(128)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (o *Order) ParseToDTOResponseCreateOrder() dto.ResponseCreateOrder {
	return dto.ResponseCreateOrder{
		ID:        o.ID,
		CanteenID: o.CanteenID,
		UserID:    o.UserID,
		Quantity:  o.Quantity,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func (o *Order) ParseToDTOResponseUpdateOrder() dto.ResponseUpdateOrder {
	return dto.ResponseUpdateOrder{
		ID:        o.ID,
		CanteenID: o.CanteenID,
		UserID:    o.UserID,
		Quantity:  o.Quantity,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func (o *Order) ParseToDTOResponseGetOrderInfo() dto.ResponseGetOrderInfo {
	return dto.ResponseGetOrderInfo{
		ID:        o.ID,
		CanteenID: o.CanteenID,
		UserID:    o.UserID,
		Quantity:  o.Quantity,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func (o *Order) ParseToDTOResponseGetOrderList() dto.ResponseGetOrderList {
	return dto.ResponseGetOrderList{
		ID:        o.ID,
		CanteenID: o.CanteenID,
		UserID:    o.UserID,
		Quantity:  o.Quantity,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
