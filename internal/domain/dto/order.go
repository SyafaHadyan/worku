// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrder struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id" validate:"required,uuid_rfc4122"`
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid_rfc4122"`
	MenuID    uuid.UUID `json:"menu_id" validate:"required,uuid_rfc4122"`
	Quantity  uint32    `json:"quantity" validate:"required,number,min=1"`
	Status    string    `json:"status"`
}

type ResponseCreateOrder struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	UserID    uuid.UUID `json:"user_id"`
	MenuID    uuid.UUID `json:"menu_id"`
	Quantity  uint32    `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateOrder struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid_rfc4122"`
	Status string    `json:"status" validate:"required,oneof=COOKING COMPLETED"`
}

type ResponseUpdateOrder struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	UserID    uuid.UUID `json:"user_id"`
	MenuID    uuid.UUID `json:"menu_id"`
	Quantity  uint32    `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetOrderInfo struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	UserID    uuid.UUID `json:"user_id"`
	MenuID    uuid.UUID `json:"menu_id"`
	Quantity  uint32    `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetOrderInfo struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	UserID    uuid.UUID `json:"user_id"`
	MenuID    uuid.UUID `json:"menu_id"`
	Quantity  uint32    `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetOrderList struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	UserID    uuid.UUID `json:"user_id"`
	MenuID    uuid.UUID `json:"menu_id"`
	Quantity  uint32    `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
