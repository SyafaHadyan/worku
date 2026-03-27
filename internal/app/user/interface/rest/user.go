// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/SyafaHadyan/worku/internal/app/user/usecase"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	googleoauth2 "github.com/SyafaHadyan/worku/internal/infra/oauth/google"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Validator   *validator.Validate
	Middleware  middleware.MiddlewareItf
	UserUseCase usecase.UserUseCaseItf
	GoogleOAuth googleoauth2.GoogleOAuthItf
	Config      *env.Env
}

func NewUserHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, userUseCase usecase.UserUseCaseItf,
	googleOAuth googleoauth2.GoogleOAuthItf, config *env.Env,
) {
	userHandler := UserHandler{
		Validator:   validator,
		Middleware:  middleware,
		UserUseCase: userUseCase,
		GoogleOAuth: googleOAuth,
		Config:      config,
	}

	routerGroup = routerGroup.Group("/users")

	routerGroup.Post("/register", userHandler.Register)
	routerGroup.Post("/login", userHandler.Login)
	routerGroup.Get("/auth/google", userHandler.GoogleLogin)
	routerGroup.Get("/auth/google/callback", userHandler.GoogleCallback)
	routerGroup.Get("/info", middleware.Authentication, userHandler.GetUserInfo)
	routerGroup.Patch("", middleware.Authentication, userHandler.UpdateUserInfo)
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var register dto.Register

	err := ctx.BodyParser(&register)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(register)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := h.UserUseCase.Register(register)
	if err != nil {
		return fiber.NewError(
			http.StatusConflict,
			"please use another email / username",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user registered",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserInfo(ctx *fiber.Ctx) error {
	var updateUserInfo dto.UpdateUserInfo

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserInfo)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserInfo)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := h.UserUseCase.UpdateUserInfo(updateUserInfo, userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user info updated",
		"payload": res,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var login dto.Login

	err := ctx.BodyParser(&login)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(login)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, token, err := h.UserUseCase.Login(login)
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"invalid email, username or password",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user authenticated",
		"token":   token,
		"payload": res,
	})
}

func (h *UserHandler) GoogleLogin(ctx *fiber.Ctx) error {
	path := h.GoogleOAuth.GoogleOAuthConfig()
	url := path.AuthCodeURL(h.GoogleOAuth.GenerateRandomState())

	return ctx.Redirect(url)
}

func (h *UserHandler) GoogleCallback(ctx *fiber.Ctx) error {
	var responseGoogleOAuth dto.ResponseGoogleOAuth

	oAuthConfig := h.GoogleOAuth.GoogleOAuthConfig()
	oAuthToken, err := oAuthConfig.Exchange(context.Background(), ctx.FormValue("code"))
	if err != nil {
		log.Println(err)
		return fiber.NewError(
			http.StatusServiceUnavailable,
			"failed to receive google's response",
		)
	}

	responseGoogleOAuth, err = h.GoogleOAuth.GetUserInfo(oAuthToken.AccessToken)
	if err != nil {
		log.Println(err)
		return fiber.NewError(
			http.StatusServiceUnavailable,
			"failed to get user info",
		)
	}

	res, token, err := h.UserUseCase.GoogleOAuth(responseGoogleOAuth)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to log in with google",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "successfully logged in with google",
		"token":   token,
		"payload": res,
	})
}

func (h *UserHandler) GetUserInfo(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserInfo(userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user info",
		"payload": res,
	})
}

func (h *UserHandler) SoftDelete(ctx *fiber.Ctx) error {
	targetUserName := ctx.Params("username")
	userIDTarget, err := h.UserUseCase.GetUserIDFromUsername(targetUserName)
	if err != nil {
		return fiber.NewError(
			http.StatusNotFound,
			"target user not found")
	}

	h.UserUseCase.SoftDelete(userIDTarget)

	return ctx.Status(http.StatusNoContent).Context().Err()
}
