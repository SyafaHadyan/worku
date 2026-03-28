// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"net/http"

	"github.com/SyafaHadyan/worku/internal/app/payment/usecase"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	Validator      *validator.Validate
	Middleware     middleware.MiddlewareItf
	PaymentUseCase usecase.PaymentUseCaseItf
}

func NewPaymentHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, paymentUseCase usecase.PaymentUseCaseItf,
) {
	paymentHandler := PaymentHandler{
		Validator:      validator,
		Middleware:     middleware,
		PaymentUseCase: paymentUseCase,
	}

	routerGroup.Post("/orders", middleware.Authentication, paymentHandler.CreateOrder)
	routerGroup.Post("/payments", middleware.Authentication, paymentHandler.CreatePayment)
	routerGroup.Post("/payments/verify", paymentHandler.VerifyPayment)
}

func (h *PaymentHandler) CreateOrder(ctx *fiber.Ctx) error {
	var createOrder dto.CreateOrder

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&createOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	createOrder.UserID = userID

	err = h.Validator.Struct(createOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := h.PaymentUseCase.CreateOrder(createOrder)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to create new order",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "order created",
		"payload": res,
	})
}

func (h *PaymentHandler) CreatePayment(ctx *fiber.Ctx) error {
	var createPayment dto.CreatePayment

	err := ctx.BodyParser(&createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := h.PaymentUseCase.CreatePayment(createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusServiceUnavailable,
			"unable to connect to 3rd party service",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "payment created",
		"payload": res,
	})
}

func (h *PaymentHandler) VerifyPayment(ctx *fiber.Ctx) error {
	var verifyPayment dto.VerifyPayment

	err := ctx.BodyParser(&verifyPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(verifyPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	err = h.PaymentUseCase.VerifyPayment(verifyPayment)
	if err == fiber.ErrBadRequest {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid payment",
		)
	} else if err == fiber.ErrPaymentRequired {
		return fiber.NewError(
			http.StatusPaymentRequired,
			"payment could not be verified",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to verify payment",
		)
	}

	return ctx.Status(http.StatusOK).Context().Err()
}
