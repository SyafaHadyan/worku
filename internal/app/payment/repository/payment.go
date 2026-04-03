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
	GetOrderInfoAfterPayment(order *entity.Order) error
	GetOrderIDFromPayment(order *entity.Order) (uuid.UUID, error)
	GetOrderList(offset *int, limit *int, userID uuid.UUID, order *[]entity.Order) error
	CreatePayment(payment *entity.Payment) error
	VerifyPayment(order *entity.Order) error
	GetUserSubscriptionExpiryDate(userSubscription *entity.UserSubscription) error
	UpdateUserPaidStatus(userSubscription *entity.UserSubscription) error
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
	return r.db.
		Create(order).
		Error
}

func (r *PaymentDB) GetOrderInfo(order *entity.Order) error {
	return r.db.
		Model(&entity.Order{}).
		Where("id = ? AND user_id = ?", order.ID, order.UserID).
		First(order).
		Error
}

func (r *PaymentDB) GetOrderInfoAfterPayment(order *entity.Order) error {
	return r.db.
		Model(&entity.Order{}).
		Where("id = ?", order.ID).
		First(order).
		Error
}

func (r *PaymentDB) GetOrderIDFromPayment(order *entity.Order) (uuid.UUID, error) {
	var payment entity.Payment

	err := r.db.
		Model(&entity.Payment{}).
		Select("order_id").
		Where("id = ?", order.ID).
		First(&payment).
		Error

	return payment.OrderID, err
}

func (r *PaymentDB) GetOrderList(offset *int, limit *int, userID uuid.UUID, order *[]entity.Order) error {
	return r.db.
		Model(&entity.Order{}).
		Where("user_id = ?", userID).
		Limit(*limit).
		Offset(*offset).
		Find(order).
		Error
}

func (r *PaymentDB) CreatePayment(payment *entity.Payment) error {
	return r.db.
		Create(payment).
		Error
}

func (r *PaymentDB) VerifyPayment(order *entity.Order) error {
	var payment entity.Payment

	r.db.
		Model(&entity.Payment{}).
		Select("order_id").
		Where("id = ?", order.ID).
		First(&payment)

	return r.db.
		Model(&entity.Order{}).
		Where("id = ?", payment.OrderID).
		Update("status", "PAID").
		Error
}

func (r *PaymentDB) GetUserSubscriptionExpiryDate(userSubscription *entity.UserSubscription) error {
	return r.db.
		Model(&entity.UserSubscription{}).
		Where("user_id = ?", userSubscription.UserID).
		First(userSubscription).
		Error
}

func (r *PaymentDB) UpdateUserPaidStatus(userSubscription *entity.UserSubscription) error {
	var err error

	if r.db.
		Model(&entity.UserSubscription{}).
		Where("user_id = ?", userSubscription.UserID).
		Update("expiry_date", userSubscription.ExpiryDate).
		RowsAffected == 0 {
		err = r.db.
			Create(userSubscription).
			Error
	}

	return err
}
