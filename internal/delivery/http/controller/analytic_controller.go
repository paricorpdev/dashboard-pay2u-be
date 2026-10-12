package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/usecase"
	"github.com/sirupsen/logrus"
)

type AnalyticController struct {
	UseCase *usecase.AnalyticUseCase
	Log     *logrus.Logger
}

func NewAnalyticController(useCase *usecase.AnalyticUseCase, log *logrus.Logger) *AnalyticController {
	return &AnalyticController{
		UseCase: useCase,
		Log:     log,
	}
}

// GetSummary returns monthly KPI summary metrics.
func (c *AnalyticController) GetSummary(ctx *fiber.Ctx) error {
	request := new(payload.AnalyticFilterRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse query params: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.GetSummary(ctx, request)
	if err != nil {
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

// GetDailyTrend returns daily transaction breakdown for line/bar charts.
func (c *AnalyticController) GetDailyTrend(ctx *fiber.Ctx) error {
	request := new(payload.AnalyticFilterRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse query params: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.GetDailyTrend(ctx, request)
	if err != nil {
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

// GetPaymentMethodDistribution returns distribution by payment method for pie/donut charts.
func (c *AnalyticController) GetPaymentMethodDistribution(ctx *fiber.Ctx) error {
	request := new(payload.AnalyticFilterRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse query params: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.GetPaymentMethodDistribution(ctx, request)
	if err != nil {
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

// GetStatusDistribution returns composition of transaction statuses for status charts.
func (c *AnalyticController) GetStatusDistribution(ctx *fiber.Ctx) error {
	request := new(payload.AnalyticFilterRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse query params: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.GetStatusDistribution(ctx, request)
	if err != nil {
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}
