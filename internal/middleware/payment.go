package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) PaidUser(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := m.userUseCase.GetUserSubscription(userID)
	if err != nil || res.ExpiryDate.Before(time.Now()) {
		return fiber.NewError(
			http.StatusPaymentRequired,
			"payment required",
		)
	}

	return ctx.Next()
}
