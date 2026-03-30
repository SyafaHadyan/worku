// Package usecase handles the logic for each user request
package usecase

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/SyafaHadyan/worku/internal/app/payment/repository"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/SyafaHadyan/worku/internal/infra/payment"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentUseCaseItf interface {
	CreateOrder(createOrder dto.CreateOrder) (dto.ResponseCreateOrder, error)
	CreateSnapPayment(createPayment dto.CreatePayment) (dto.ResponseCreateMidtransOrder, error)
	GetOrderInfo(getOrderInfo dto.GetOrderInfo) (dto.ResponseGetOrderInfo, error)
	GetOrderList(offset int, limit int, userID uuid.UUID) ([]dto.ResponseGetOrderList, error)
	VerifyPayment(verifyPayment dto.VerifyPayment) error
}

type PaymentUseCase struct {
	paymentRepo repository.PaymentDBItf
	payment     payment.PaymentItf
	env         *env.Env
}

func NewPaymentuseCase(
	paymentRepo repository.PaymentDBItf, payment payment.PaymentItf,
	env *env.Env,
) PaymentUseCaseItf {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
		payment:     payment,
		env:         env,
	}
}

func (u *PaymentUseCase) CreateOrder(createOrder dto.CreateOrder) (dto.ResponseCreateOrder, error) {
	order := entity.Order{
		ID:           uuid.New(),
		UserID:       createOrder.UserID,
		DurationDays: createOrder.DurationDays,
		Status:       "UNPAID",
	}

	err := u.paymentRepo.CreateOrder(&order)

	return order.ParseToDTOResponseCreateOrder(), err
}

func (u *PaymentUseCase) CreateSnapPayment(createPayment dto.CreatePayment) (dto.ResponseCreateMidtransOrder, error) {
	paymentID := uuid.New()
	order := entity.Order{
		ID:     createPayment.OrderID,
		UserID: createPayment.UserID,
	}

	err := u.paymentRepo.GetOrderInfo(&order)
	if err != nil || order.Status == "PAID" {
		return dto.ResponseCreateMidtransOrder{}, gorm.ErrRecordNotFound
	}

	createMidtransSnapOrder := dto.CreateMidtransSnapOrder{
		TransactionDetails: dto.TransactionDetails{
			OrderID: paymentID.String(),
		},
		Interval: order.DurationDays,
	}

	res, err := u.payment.CreateSnapPayment(createMidtransSnapOrder)

	payment := entity.Payment{
		ID:          paymentID,
		OrderID:     createPayment.OrderID,
		Token:       res.Token,
		RedirectURL: res.RedirectURL,
	}

	err = u.paymentRepo.CreatePayment(&payment)

	return payment.ParseToDTOResponseCreateMidtransOrder(), err
}

func (u *PaymentUseCase) CreateCoreAPIPayment(createPayment dto.CreatePayment) (dto.ResponseCreateMidtransOrder, error) {
	paymentID := uuid.New()
	order := entity.Order{
		ID:     createPayment.OrderID,
		UserID: createPayment.UserID,
	}

	err := u.paymentRepo.GetOrderInfo(&order)
	if err != nil || order.Status == "PAID" {
		return dto.ResponseCreateMidtransOrder{}, gorm.ErrRecordNotFound
	}

	createMidtransCoreAPIOrder := dto.CreateMidtransCoreAPIOrder{
		TransactionDetails: dto.TransactionDetails{
			OrderID: paymentID.String(),
		},
		Interval: order.DurationDays,
	}

	res, err := u.payment.CreateCoreAPIPayment(createMidtransCoreAPIOrder)

	payment := entity.Payment{
		ID:          paymentID,
		OrderID:     createPayment.OrderID,
		RedirectURL: res.RedirectURL,
	}

	err = u.paymentRepo.CreatePayment(&payment)

	return payment.ParseToDTOResponseCreateMidtransOrder(), err
}

func (u *PaymentUseCase) GetOrderInfo(getOrderInfo dto.GetOrderInfo) (dto.ResponseGetOrderInfo, error) {
	order := entity.Order{
		ID:     getOrderInfo.ID,
		UserID: getOrderInfo.UserID,
	}

	err := u.paymentRepo.GetOrderInfo(&order)

	return order.ParseToDTOResponseGetOrderInfo(), err
}

func (u *PaymentUseCase) GetOrderList(offset int, limit int, userID uuid.UUID) ([]dto.ResponseGetOrderList, error) {
	var order []entity.Order

	offset = offset * limit

	err := u.paymentRepo.GetOrderList(&offset, &limit, userID, &order)
	if err != nil {
		return nil, err
	}

	if len(order) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	orderList := make([]dto.ResponseGetOrderList, len(order))

	for i, courseItem := range order {
		orderList[i] = courseItem.ParseToDTOResponseGetOrderList()
	}

	return orderList, nil
}

func (c *PaymentUseCase) VerifyPayment(verifyPayment dto.VerifyPayment) error {
	orderID, _ := uuid.Parse(verifyPayment.OrderID)
	transactionStatus := verifyPayment.TransactionStatus

	signatureKey := fmt.Sprintf(
		"%s%s%s%s",
		verifyPayment.OrderID,
		verifyPayment.StatusCode,
		verifyPayment.GrossAmount,
		c.env.MidtransServerKey,
	)

	hash := sha512.New()
	hash.Write([]byte(signatureKey))
	hashedData := hash.Sum(nil)
	hexHash := hex.EncodeToString(hashedData)

	if hexHash != verifyPayment.SignatureKey {
		return fiber.ErrBadRequest
	}

	if transactionStatus != "capture" && transactionStatus != "settlement" {
		return fiber.ErrPaymentRequired
	}
	order := entity.Order{
		ID: orderID,
	}

	err := c.paymentRepo.VerifyPayment(&order)

	return err
}
