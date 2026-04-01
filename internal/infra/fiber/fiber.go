// Package fiber sets configuration for the Fiber framework
package fiber

import (
	"time"

	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Fiber struct {
	Fiber  *fiber.App
	Router fiber.Router
}

func New(env *env.Env) *Fiber {
	app := fiber.New(
		fiber.Config{
			Prefork:   false,
			BodyLimit: env.BodyLimit * 1024 * 1024,
		},
	)

	app.Use(
		idempotency.New(),
		cors.New(
			cors.Config{
				AllowHeaders: "*",
				AllowOrigins: "*",
				AllowMethods: "*",
			}),
		limiter.New(
			limiter.Config{
				Max:               env.LimiterMax,
				Expiration:        time.Duration(env.LimiterExpirationMinutes) * 60,
				LimiterMiddleware: limiter.SlidingWindow{},
			}),
		logger.New(
			logger.Config{
				Format: "${ip} - - [${time}] ${method} ${url} ${protocol} ${status} ${bytesSent} ${referer} ${ua}\n",
			},
		),
	)

	v1 := app.Group("/api/v1")

	Fiber := Fiber{
		Fiber:  app,
		Router: v1,
	}

	return &Fiber
}
