// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePayment struct {
	ID          uuid.UUID      `json:"id"`
	OrderID     uuid.UUID      `json:"order_id" validate:"required,uuid_rfc4122"`
	UserID      uuid.UUID      `json:"user_id"`
	Price       uint32         `json:"price"`
	RedirectURL string         `json:"redirect_url"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CreateMidtransOrder struct {
	TransactionDetails TransactionDetails
	CustomerDetail     CustomerDetail
}

type TransactionDetails struct {
	OrderID     string `json:"order_id" validate:"required,number,min=1"`
	GrossAmount uint32 `json:"gross_amount" validate:"required,number,min=1"`
}

type CustomerDetail struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1"`
	LastName  string `json:"last_name" validate:"omitempty,min=1"`
	Email     string `json:"email" validate:"omitempty,email"`
}

type ResponseMidtransOrder struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type VerifyPayment struct {
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
}
