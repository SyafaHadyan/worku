// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"net/http"
	"strconv"

	"github.com/SyafaHadyan/worku/internal/app/payment/usecase"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	routerGroup.Post("/payments/snap", middleware.Authentication, paymentHandler.CreateSnapPayment)
	routerGroup.Post("/payments/coreapi", middleware.Authentication, paymentHandler.CreateCoreAPIPayment)
	routerGroup.Get("/orders/:id", middleware.Authentication, paymentHandler.GetOrderInfo)
	routerGroup.Get("orders/list/:page/:limit", middleware.Authentication, paymentHandler.GetOrderList)
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

func (h *PaymentHandler) CreateSnapPayment(ctx *fiber.Ctx) error {
	var createPayment dto.CreatePayment

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	createPayment.UserID = userID

	err = h.Validator.Struct(createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := h.PaymentUseCase.CreateSnapPayment(createPayment)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"order not found / already paid",
		)
	} else if err != nil {
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

func (h *PaymentHandler) CreateCoreAPIPayment(ctx *fiber.Ctx) error {
	var createPayment dto.CreatePayment

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	createPayment.UserID = userID

	err = h.Validator.Struct(createPayment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := h.PaymentUseCase.CreateSnapPayment(createPayment)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"order not found / already paid",
		)
	} else if err != nil {
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

func (h *PaymentHandler) GetOrderInfo(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	orderID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid order id",
		)
	}

	getOrderInfo := dto.GetOrderInfo{
		ID:     orderID,
		UserID: userID,
	}

	res, err := h.PaymentUseCase.GetOrderInfo(getOrderInfo)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"order not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get order info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully get order info",
		"payload": res,
	})
}

func (h *PaymentHandler) GetOrderList(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	offset, err := strconv.Atoi(ctx.Params("page", "0"))
	if err != nil || offset < 0 {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid page",
		)
	}

	limit, err := strconv.Atoi(ctx.Params("limit", "8"))
	if err != nil || limit <= 0 {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid limit",
		)
	}

	res, err := h.PaymentUseCase.GetOrderList(offset, limit, userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get order list",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved order list",
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
	if err == fiber.ErrBadRequest || err == gorm.ErrRecordNotFound {
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
