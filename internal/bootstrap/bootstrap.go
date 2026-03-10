// Package bootstrap starts the backend service, sets the backend configuration, and connects to external services
package bootstrap

import (
	"log"
	"time"

	coursehandler "github.com/SyafaHadyan/worku/internal/app/course/interface/rest"
	courserepository "github.com/SyafaHadyan/worku/internal/app/course/repository"
	courseusecase "github.com/SyafaHadyan/worku/internal/app/course/usecase"
	userhandler "github.com/SyafaHadyan/worku/internal/app/user/interface/rest"
	userrepository "github.com/SyafaHadyan/worku/internal/app/user/repository"
	userusecase "github.com/SyafaHadyan/worku/internal/app/user/usecase"
	"github.com/SyafaHadyan/worku/internal/infra/db"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	fiberapp "github.com/SyafaHadyan/worku/internal/infra/fiber"
	"github.com/SyafaHadyan/worku/internal/infra/jwt"
	"github.com/SyafaHadyan/worku/internal/infra/payment"
	"github.com/SyafaHadyan/worku/internal/infra/redis"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
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

	database := db.New(config)

	redis := redis.New(config)

	jwt := jwt.New(config)

	_ = payment.New(config)

	app := fiberapp.New(config)

	middleware := middleware.NewMiddleware(*jwt)

	userRepository := userrepository.NewUserDB(database)
	courseRepository := courserepository.NewCourseDB(database)

	userUseCase := userusecase.NewUserUseCase(userRepository, jwt, redis)
	courseUseCase := courseusecase.NewCourseUseCase(courseRepository, redis)

	userhandler.NewUserHandler(app.Router, validator, middleware, userUseCase, config)
	coursehandler.NewCourseHandler(app.Router, validator, middleware, courseUseCase)

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
