// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"net/http"
	"strconv"

	"github.com/SyafaHadyan/worku/internal/app/course/usecase"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseHandler struct {
	Validator     *validator.Validate
	Middleware    middleware.MiddlewareItf
	CourseUseCase usecase.CourseUseCaseItf
}

func NewCourseHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, courseUseCase usecase.CourseUseCaseItf,
) {
	courseHandler := CourseHandler{
		Validator:     validator,
		Middleware:    middleware,
		CourseUseCase: courseUseCase,
	}

	routerGroup = routerGroup.Group("/courses")

	routerGroup.Get("/list/:page/:limit", middleware.Authentication, courseHandler.GetCourseList)
	routerGroup.Get("/:id", middleware.Authentication, courseHandler.GetCourseInfo)
	routerGroup.Get("/search/:query", middleware.Authentication, courseHandler.SearchCourse)
}

func (h *CourseHandler) GetCourseList(ctx *fiber.Ctx) error {
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

	res, err := h.CourseUseCase.GetCourseList(offset, limit)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get course list",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved course list",
		"payload": res,
	})
}

func (h *CourseHandler) GetCourseInfo(ctx *fiber.Ctx) error {
	courseID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid course id",
		)
	}

	res, err := h.CourseUseCase.GetCourseInfo(courseID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"course not found",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get course info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved course info",
		"payload": res,
	})
}

func (h *CourseHandler) SearchCourse(ctx *fiber.Ctx) error {
	query := ctx.Params("query")

	res, err := h.CourseUseCase.SearchCourse(query)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"search query doesn't match any course",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to search course",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"messsage": "retrieved course list",
		"payload":  res,
	})
}
