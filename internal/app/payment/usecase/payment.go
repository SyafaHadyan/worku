// Package usecase handles the logic for each user request
package usecase

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/SyafaHadyan/worku/internal/app/payment/repository"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/domain/entity"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/SyafaHadyan/worku/internal/infra/payment"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentUseCaseItf interface {
	CreateOrder(createOrder dto.CreateOrder) (dto.ResponseCreateOrder, error)
	CreatePayment(createPayment dto.CreatePayment) (dto.ResponseCreateMidtransOrder, error)
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

func (u *PaymentUseCase) CreatePayment(createPayment dto.CreatePayment) (dto.ResponseCreateMidtransOrder, error) {
	paymentID := uuid.New()
	order := entity.Order{
		ID: createPayment.OrderID,
	}

	err := u.paymentRepo.GetOrderInfo(&order)

	createMidtransOrder := dto.CreateMidtransOrder{
		TransactionDetails: dto.TransactionDetails{
			OrderID: paymentID.String(),
		},
		Interval: order.DurationDays,
	}
	res, err := u.payment.CreatePayment(createMidtransOrder)

	payment := entity.Payment{
		ID:          paymentID,
		OrderID:     createPayment.OrderID,
		Token:       res.Token,
		RedirectURL: res.RedirectURL,
	}

	err = u.paymentRepo.CreatePayment(&payment)

	return payment.ParseToDTOResponseCreateMidtransOrder(), err
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

	log.Println(order)

	err := c.paymentRepo.VerifyPayment(&order)

	return err
}
