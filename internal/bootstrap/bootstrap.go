// Package bootstrap starts the backend service, sets the backend configuration, and connects to external services
package bootstrap

import (
	"log"
	"time"

	aihandler "github.com/SyafaHadyan/worku/internal/app/ai/interface/rest"
	airepository "github.com/SyafaHadyan/worku/internal/app/ai/repository"
	aiusecase "github.com/SyafaHadyan/worku/internal/app/ai/usecase"
	coursehandler "github.com/SyafaHadyan/worku/internal/app/course/interface/rest"
	courserepository "github.com/SyafaHadyan/worku/internal/app/course/repository"
	courseusecase "github.com/SyafaHadyan/worku/internal/app/course/usecase"
	paymenthandler "github.com/SyafaHadyan/worku/internal/app/payment/interface/rest"
	paymentrepository "github.com/SyafaHadyan/worku/internal/app/payment/repository"
	paymentusecase "github.com/SyafaHadyan/worku/internal/app/payment/usecase"
	userhandler "github.com/SyafaHadyan/worku/internal/app/user/interface/rest"
	userrepository "github.com/SyafaHadyan/worku/internal/app/user/repository"
	userusecase "github.com/SyafaHadyan/worku/internal/app/user/usecase"
	"github.com/SyafaHadyan/worku/internal/infra/ai"
	"github.com/SyafaHadyan/worku/internal/infra/db"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	fiberapp "github.com/SyafaHadyan/worku/internal/infra/fiber"
	"github.com/SyafaHadyan/worku/internal/infra/jwt"
	googleoauth2 "github.com/SyafaHadyan/worku/internal/infra/oauth/google"
	"github.com/SyafaHadyan/worku/internal/infra/payment"
	"github.com/SyafaHadyan/worku/internal/infra/redis"
	"github.com/SyafaHadyan/worku/internal/infra/s3"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type Bootstrap struct {
	App       *fiberapp.Fiber
	Config    *env.Env
	Validator *validator.Validate
	Database  *gorm.DB
	Redis     *redis.Redis
	JWT       *jwt.JWT
}

func Start() *Bootstrap {
	log.Println("starting app")
	startTime := time.Now()

	config := env.New()

	validator := validator.New()

	decoder := schema.NewDecoder()

	database := db.New(config)

	redis := redis.New(config)

	ai := ai.New(config)

	s3 := s3.New(config)

	jwt := jwt.New(config)

	googleoauth2 := googleoauth2.New(config)

	payment := payment.New(config)

	app := fiberapp.New(config)

	userRepository := userrepository.NewUserDB(database)
	courseRepository := courserepository.NewCourseDB(database)
	aiRepository := airepository.NewAIDB(database)
	paymentRepository := paymentrepository.NewPaymentDB(database)

	userUseCase := userusecase.NewUserUseCase(userRepository, jwt, redis, s3, config)
	courseUseCase := courseusecase.NewCourseUseCase(courseRepository, redis)
	aiUseCase := aiusecase.NewAIUseCase(aiRepository, ai)
	paymentUseCase := paymentusecase.NewPaymentuseCase(paymentRepository, payment, config)

	middleware := middleware.NewMiddleware(*jwt, userUseCase)

	userhandler.NewUserHandler(app.Router, validator, middleware, userUseCase, googleoauth2, config)
	coursehandler.NewCourseHandler(app.Router, validator, middleware, courseUseCase)
	aihandler.NewAIHandler(app.Router, validator, decoder, middleware, aiUseCase)
	paymenthandler.NewPaymentHandler(app.Router, validator, middleware, paymentUseCase)

	log.Printf("startup time: %v", time.Since(startTime))

	Bootstrap := Bootstrap{
		App:       app,
		Config:    config,
		Validator: validator,
		Database:  database,
		Redis:     redis,
		JWT:       jwt,
	}

	return &Bootstrap
}
