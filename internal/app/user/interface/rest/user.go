// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/SyafaHadyan/worku/internal/app/user/usecase"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	googleoauth2 "github.com/SyafaHadyan/worku/internal/infra/oauth/google"
	linkedinoauth2 "github.com/SyafaHadyan/worku/internal/infra/oauth/linkedin"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	Validator     *validator.Validate
	Middleware    middleware.MiddlewareItf
	UserUseCase   usecase.UserUseCaseItf
	GoogleOAuth   googleoauth2.GoogleOAuthItf
	LinkedInOAuth linkedinoauth2.LinkedInOAuthItf
	Config        *env.Env
}

func NewUserHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, userUseCase usecase.UserUseCaseItf,
	googleOAuth googleoauth2.GoogleOAuthItf, linkedinOAuth linkedinoauth2.LinkedInOAuthItf,
	config *env.Env,
) {
	userHandler := UserHandler{
		Validator:     validator,
		Middleware:    middleware,
		UserUseCase:   userUseCase,
		GoogleOAuth:   googleOAuth,
		LinkedInOAuth: linkedinOAuth,
		Config:        config,
	}

	routerGroup = routerGroup.Group("/users")

	routerGroup.Post("/register", userHandler.Register)
	routerGroup.Post("/login", userHandler.Login)
	routerGroup.Patch("", middleware.Authentication, userHandler.UpdateUserInfo)
	routerGroup.Patch("/info/detail", middleware.Authentication, userHandler.UpdateUserDetail)
	routerGroup.Patch("/info/contact", middleware.Authentication, userHandler.UpdateUserContact)
	routerGroup.Patch("/info/education", middleware.Authentication, userHandler.UpdateUserEducation)
	routerGroup.Post("/info/language", middleware.Authentication, userHandler.AddUserLanguage)
	routerGroup.Patch("/info/employment", middleware.Authentication, userHandler.UpdateUserEmployment)
	routerGroup.Patch("/info/seniority", middleware.Authentication, userHandler.UpdateUserSeniority)
	routerGroup.Patch("/info/workexperience", middleware.Authentication, userHandler.UpdateUserWorkExperience)
	routerGroup.Post("/info/hardskill", middleware.Authentication, userHandler.AddUserHardSkill)
	routerGroup.Post("/info/softskill", middleware.Authentication, userHandler.AddUserSoftSkill)
	routerGroup.Post("/info/tool", middleware.Authentication, userHandler.AddUserTools)
	routerGroup.Patch("/info/link", middleware.Authentication, userHandler.UpdateUserLink)
	routerGroup.Get("/auth/google", userHandler.GoogleLogin)
	routerGroup.Get("/auth/google/callback", userHandler.GoogleCallback)
	routerGroup.Get("/auth/linkedin", userHandler.LinkedInLogin)
	routerGroup.Get("/auth/linkedin/callback", userHandler.LinkedInCallback)
	routerGroup.Post("/profile/upload", middleware.Authentication, userHandler.UploadProfilePicture)
	routerGroup.Get("/info", middleware.Authentication, userHandler.GetUserInfo)
	routerGroup.Get("/info/detail", middleware.Authentication, userHandler.GetUserDetail)
	routerGroup.Get("/info/contact", middleware.Authentication, userHandler.GetUserContact)
	routerGroup.Get("/info/education", middleware.Authentication, userHandler.GetUserEducation)
	routerGroup.Get("/info/language", middleware.Authentication, userHandler.GetUserLanguage)
	routerGroup.Get("/info/employment", middleware.Authentication, userHandler.GetUserEmployment)
	routerGroup.Get("/info/seniority", middleware.Authentication, userHandler.GetUserSeniority)
	routerGroup.Get("/info/workexperience", middleware.Authentication, userHandler.GetUserWorkExperience)
	routerGroup.Get("/info/hardskill", middleware.Authentication, userHandler.GetUserHardSkill)
	routerGroup.Get("/info/softskill", middleware.Authentication, userHandler.GetUserSoftSkill)
	routerGroup.Get("/info/tool", middleware.Authentication, userHandler.GetUserTools)
	routerGroup.Get("/info/link", middleware.Authentication, userHandler.GetUserLink)
	routerGroup.Get("/info/subscription", middleware.Authentication, userHandler.GetUserSubscription)
	routerGroup.Delete("/info/language/:language", middleware.Authentication, userHandler.DeleteUserLanguage)
	routerGroup.Delete("/info/hardskill/:hardskill", middleware.Authentication, userHandler.DeleteUserHardSkill)
	routerGroup.Delete("/info/softskill/:softskill", middleware.Authentication, userHandler.DeleteUserSoftSkill)
	routerGroup.Delete("/info/tool/:tool", middleware.Authentication, userHandler.DeleteUserTools)
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

	updateUserInfo.ID = userID

	res, err := h.UserUseCase.UpdateUserInfo(updateUserInfo)
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

func (h *UserHandler) UpdateUserDetail(ctx *fiber.Ctx) error {
	var updateUserDetail dto.UpdateUserDetail

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserDetail)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserDetail)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserDetail.UserID = userID

	res, err := h.UserUseCase.UpdateUserDetail(updateUserDetail)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user detail",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user detail updated",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserContact(ctx *fiber.Ctx) error {
	var updateUserContact dto.UpdateUserContact

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserContact)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserContact)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserContact.UserID = userID

	res, err := h.UserUseCase.UpdateUserContact(updateUserContact)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user contact",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user contact updated",
		"payload": res,
	})
}

func (h *UserHandler) AddUserLanguage(ctx *fiber.Ctx) error {
	var addUserLanguage dto.AddUserLanguage

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&addUserLanguage)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(addUserLanguage)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	addUserLanguage.UserID = userID

	res, err := h.UserUseCase.AddUserLanguage(addUserLanguage)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to add user language",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "added user language",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserEducation(ctx *fiber.Ctx) error {
	var updateUserEducation dto.UpdateUserEducation

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserEducation)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserEducation)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserEducation.UserID = userID

	res, err := h.UserUseCase.UpdateUserEducation(updateUserEducation)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user education",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user education updated",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserEmployment(ctx *fiber.Ctx) error {
	var updateUserEmployment dto.UpdateUserEmployment

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserEmployment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserEmployment)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserEmployment.UserID = userID

	res, err := h.UserUseCase.UpdateUserEmployment(updateUserEmployment)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user employment",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user employment updated",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserSeniority(ctx *fiber.Ctx) error {
	var updateUserSeniority dto.UpdateUserSeniority

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserSeniority)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserSeniority)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserSeniority.UserID = userID

	res, err := h.UserUseCase.UpdateUserSeniority(updateUserSeniority)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user seniority",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user seniority updated",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserWorkExperience(ctx *fiber.Ctx) error {
	var updateUserWorkExperience dto.UpdateUserWorkExperience

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserWorkExperience)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserWorkExperience)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserWorkExperience.UserID = userID

	res, err := h.UserUseCase.UpdateUserWorkExperience(updateUserWorkExperience)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user work experience",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user work experience updated",
		"payload": res,
	})
}

func (h *UserHandler) AddUserHardSkill(ctx *fiber.Ctx) error {
	var addUserHardSkill dto.AddUserHardSkill

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&addUserHardSkill)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(addUserHardSkill)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	addUserHardSkill.UserID = userID

	res, err := h.UserUseCase.AddUserHardSkill(addUserHardSkill)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to add user hard skill",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "added user hard skill",
		"payload": res,
	})
}

func (h *UserHandler) AddUserSoftSkill(ctx *fiber.Ctx) error {
	var addUserSoftSkill dto.AddUserSoftSkill

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&addUserSoftSkill)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(addUserSoftSkill)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	addUserSoftSkill.UserID = userID

	res, err := h.UserUseCase.AddUserSoftSkill(addUserSoftSkill)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to add user soft skill",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "added user soft skill",
		"payload": res,
	})
}

func (h *UserHandler) AddUserTools(ctx *fiber.Ctx) error {
	var addUserTools dto.AddUserTools

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&addUserTools)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(addUserTools)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	addUserTools.UserID = userID

	res, err := h.UserUseCase.AddUserTools(addUserTools)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to add user tools",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "added user tools",
		"payload": res,
	})
}

func (h *UserHandler) UpdateUserLink(ctx *fiber.Ctx) error {
	var updateUserLink dto.UpdateUserLink

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&updateUserLink)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = h.Validator.Struct(updateUserLink)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updateUserLink.UserID = userID

	res, err := h.UserUseCase.UpdateUserLink(updateUserLink)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user link",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user link updated",
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

func (h *UserHandler) LinkedInLogin(ctx *fiber.Ctx) error {
	path := h.LinkedInOAuth.LinkedInOAuthConfig()
	url := path.AuthCodeURL(h.LinkedInOAuth.GenerateRandomState())

	return ctx.Redirect(url)
}

func (h *UserHandler) LinkedInCallback(ctx *fiber.Ctx) error {
	var responseLinkedInOAuth dto.ResponseLinkedInOAuth

	oAuthConfig := h.LinkedInOAuth.LinkedInOAuthConfig()
	oAuthToken, err := oAuthConfig.Exchange(context.Background(), ctx.FormValue("code"))
	if err != nil {
		log.Println(err)
		return fiber.NewError(
			http.StatusServiceUnavailable,
			"failed to receive linkedin's response",
		)
	}

	responseLinkedInOAuth, err = h.LinkedInOAuth.GetUserInfo(oAuthToken.AccessToken)
	if err != nil {
		log.Println(err)
		return fiber.NewError(
			http.StatusServiceUnavailable,
			"failed to get user info",
		)
	}

	res, token, err := h.UserUseCase.LinkedInOAuth(responseLinkedInOAuth)
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

func (h *UserHandler) UploadProfilePicture(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	file, err := ctx.FormFile("picture")
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to upload voice",
		)
	}

	res, err := h.UserUseCase.UploadProfilePicture(userID, *file)
	if err != nil {
		return fiber.NewError(

			http.StatusServiceUnavailable,
			"failed to connect to 3rd party service",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "successfully uploaded profile picture",
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

func (h *UserHandler) GetUserDetail(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserDetail(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user detail",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user detail",
		"payload": res,
	})
}

func (h *UserHandler) GetUserContact(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserContact(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user contact",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user contact",
		"payload": res,
	})
}

func (h *UserHandler) GetUserEducation(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserEducation(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user education",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user education",
		"payload": res,
	})
}

func (h *UserHandler) GetUserLanguage(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserLanguage(userID)
	if err == gorm.ErrRecordNotFound || len(res) == 0 {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user language",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user language",
		"payload": res,
	})
}

func (h *UserHandler) GetUserEmployment(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserEmployment(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user employment",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user employment",
		"payload": res,
	})
}

func (h *UserHandler) GetUserSeniority(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserSeniority(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user seniority",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user seniority",
		"payload": res,
	})
}

func (h *UserHandler) GetUserWorkExperience(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserWorkExperience(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user work experience",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user work experience",
		"payload": res,
	})
}

func (h *UserHandler) GetUserHardSkill(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserHardSkill(userID)
	if err == gorm.ErrRecordNotFound || len(res) == 0 {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user hard skill",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user hard skill",
		"payload": res,
	})
}

func (h *UserHandler) GetUserSoftSkill(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserSoftSkill(userID)
	if err == gorm.ErrRecordNotFound || len(res) == 0 {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user soft skill",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user soft skill",
		"payload": res,
	})
}

func (h *UserHandler) GetUserTools(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserTools(userID)
	if err == gorm.ErrRecordNotFound || len(res) == 0 {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user tools",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user tools",
		"payload": res,
	})
}

func (h *UserHandler) GetUserLink(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserLink(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user link",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user link",
		"payload": res,
	})
}

func (h *UserHandler) GetUserSubscription(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := h.UserUseCase.GetUserSubscription(userID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user subscription",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user subscription",
		"payload": res,
	})
}

func (h *UserHandler) DeleteUserLanguage(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	targetLanguage := ctx.Params("language")
	targetLanguage = strings.Replace(targetLanguage, "%20", " ", -1)

	deleteUserLanguage := dto.DeleteUserLanguage{
		UserID:         userID,
		LanguageSpoken: targetLanguage,
	}

	err = h.UserUseCase.DeleteUserLanguage(deleteUserLanguage)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"target language not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete target language",
		)
	}

	return ctx.Status(http.StatusNoContent).Context().Err()
}

func (h *UserHandler) DeleteUserHardSkill(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	targetHardSkill := ctx.Params("hardskill")
	targetHardSkill = strings.Replace(targetHardSkill, "%20", " ", -1)

	deleteUserHardSkill := dto.DeleteUserHardSkill{
		UserID:    userID,
		HardSkill: targetHardSkill,
	}

	err = h.UserUseCase.DeleteUserHardSkill(deleteUserHardSkill)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"target hard skill not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete target hard skill",
		)
	}

	return ctx.Status(http.StatusNoContent).Context().Err()
}

func (h *UserHandler) DeleteUserSoftSkill(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	targetSoftSkill := ctx.Params("softskill")
	targetSoftSkill = strings.Replace(targetSoftSkill, "%20", " ", -1)

	deleteUserSoftSkill := dto.DeleteUserSoftSkill{
		UserID:    userID,
		SoftSkill: targetSoftSkill,
	}

	err = h.UserUseCase.DeleteUserSoftSkill(deleteUserSoftSkill)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"target soft skill not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete target soft skill",
		)
	}

	return ctx.Status(http.StatusNoContent).Context().Err()
}

func (h *UserHandler) DeleteUserTools(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	targetTools := ctx.Params("tool")
	targetTools = strings.Replace(targetTools, "%20", " ", -1)

	deleteUserTools := dto.DeleteUserTools{
		UserID: userID,
		Tools:  targetTools,
	}

	err = h.UserUseCase.DeleteUserTools(deleteUserTools)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"target tools not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete target tools",
		)
	}

	return ctx.Status(http.StatusNoContent).Context().Err()
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
