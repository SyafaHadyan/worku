// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentDBItf interface {
	CreateOrder(order *entity.Order) error
	GetOrderInfo(order *entity.Order) error
	GetOrderList(offset *int, limit *int, userID uuid.UUID, order *[]entity.Order) error
	CreatePayment(payment *entity.Payment) error
	VerifyPayment(order *entity.Order) error
}

type PaymentDB struct {
	db *gorm.DB
}

func NewPaymentDB(db *gorm.DB) PaymentDBItf {
	return &PaymentDB{
		db: db,
	}
}

func (r *PaymentDB) CreateOrder(order *entity.Order) error {
	return r.db.Debug().
		Create(order).
		Error
}

func (r *PaymentDB) GetOrderInfo(order *entity.Order) error {
	return r.db.Debug().
		Model(&entity.Order{}).
		Where("id = ? AND user_id = ?", order.ID, order.UserID).
		First(order).
		Error
}

func (r *PaymentDB) GetOrderList(offset *int, limit *int, userID uuid.UUID, order *[]entity.Order) error {
	return r.db.Debug().
		Model(&entity.Order{}).
		Where("user_id = ?", userID).
		Limit(*limit).
		Offset(*offset).
		Find(order).
		Error
}

func (r *PaymentDB) CreatePayment(payment *entity.Payment) error {
	return r.db.Debug().
		Create(payment).
		Error
}

func (r *PaymentDB) VerifyPayment(order *entity.Order) error {
	var payment entity.Payment

	r.db.Debug().
		Model(&entity.Payment{}).
		Select("order_id").
		Where("id = ?", order.ID).
		First(&payment)

	return r.db.Debug().
		Model(&entity.Order{}).
		Where("id = ?", payment.OrderID).
		Update("status", "PAID").
		Error
}
