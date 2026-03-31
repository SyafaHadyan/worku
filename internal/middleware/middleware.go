// Package middleware provides interface to some basic middleware functions
package middleware

import (
	"github.com/SyafaHadyan/worku/internal/app/user/usecase"
	"github.com/SyafaHadyan/worku/internal/infra/jwt"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareItf interface {
	Authentication(ctx *fiber.Ctx) error
	PaidUser(ctx *fiber.Ctx) error
}

type Middleware struct {
	jwt         jwt.JWT
	userUseCase usecase.UserUseCaseItf
}

func NewMiddleware(jwt jwt.JWT, userUseCase usecase.UserUseCaseItf) MiddlewareItf {
	return &Middleware{
		jwt:         jwt,
		userUseCase: userUseCase,
	}
}
