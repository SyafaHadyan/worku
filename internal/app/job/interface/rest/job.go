// Package rest receive request from user and return appropriate response based on package usecase
package rest

import (
	"net/http"
	"strconv"

	"github.com/SyafaHadyan/worku/internal/app/job/usecase"
	"github.com/SyafaHadyan/worku/internal/constants"
	"github.com/SyafaHadyan/worku/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobHandler struct {
	Validator  *validator.Validate
	Middleware middleware.MiddlewareItf
	JobUseCase usecase.JobUseCaseItf
}

func NewJobHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, jobUseCase usecase.JobUseCaseItf,
) {
	jobHandler := JobHandler{
		Validator:  validator,
		Middleware: middleware,
		JobUseCase: jobUseCase,
	}

	routerGroup.Get("/job/:jobid", jobHandler.GetJobInfo)
	routerGroup.Get("/job/:page/:limit", jobHandler.GetJobList)
	routerGroup.Get("/job/search/:page/:limit/:query", jobHandler.SearchJob)
	routerGroup.Get("/company/:companyid", jobHandler.GetCompanyInfo)
}

func (h *JobHandler) GetJobInfo(ctx *fiber.Ctx) error {
	jobID, err := uuid.Parse(ctx.Params("jobid"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid job id",
		)
	}

	res, err := h.JobUseCase.GetJobInfo(jobID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"invalid job id",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get job info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved job info",
		"payload": res,
	})
}

func (h *JobHandler) GetJobList(ctx *fiber.Ctx) error {
	offset, err := strconv.Atoi(ctx.Params("page", string(constants.DefaultPage)))
	if err != nil || offset < 0 {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid page",
		)
	}

	limit, err := strconv.Atoi(ctx.Params("limit", string(constants.DefaultPage)))
	if err != nil || limit <= 0 {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid limit",
		)
	}

	res, err := h.JobUseCase.GetJobList(offset, limit)
	if err == gorm.ErrRecordNotFound || len(res) == 0 {
		return fiber.NewError(
			http.StatusNoContent)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get job list",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved job list",
		"payload": res,
	})
}

func (h *JobHandler) SearchJob(ctx *fiber.Ctx) error {
	offset, err := strconv.Atoi(ctx.Params("page", string(constants.DefaultPage)))
	if err != nil || offset < 0 {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid page",
		)
	}

	limit, err := strconv.Atoi(ctx.Params("limit", string(constants.DefaultLimit)))
	if err != nil || limit <= 0 {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid limit",
		)
	}

	query := ctx.Params("query")

	res, err := h.JobUseCase.SearchJob(offset, limit, query)
	if err == gorm.ErrRecordNotFound || len(res) == 0 {
		return fiber.NewError(
			http.StatusNotFound,
			"search query doesn't match any job",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to search job",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"messsage": "retrieved job list",
		"payload":  res,
	})
}

func (h *JobHandler) GetCompanyInfo(ctx *fiber.Ctx) error {
	companyID, err := uuid.Parse(ctx.Params("companyid"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid company id",
		)
	}

	res, err := h.JobUseCase.GetCompanyInfo(companyID)
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			http.StatusNotFound,
			"invalid company id",
		)
	} else if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get company info",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved company info",
		"payload": res,
	})
}
