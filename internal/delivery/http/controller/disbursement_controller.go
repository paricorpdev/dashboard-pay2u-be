package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/usecase"
	"github.com/sirupsen/logrus"
)

type DisbursementController struct {
	UseCase *usecase.DisbursementUseCase
	Log     *logrus.Logger
}

func NewDisbursementController(useCase *usecase.DisbursementUseCase, log *logrus.Logger) *DisbursementController {
	return &DisbursementController{UseCase: useCase, Log: log}
}

func (c *DisbursementController) FindAll(ctx *fiber.Ctx) error {
	request := new(payload.QueryDisbursementRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse disbursement query params: %+v", err)
		return c.badRequest(ctx)
	}
	result, err := c.UseCase.FindAll(ctx, request)
	if err != nil {
		return err
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: result,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *DisbursementController) FindByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return c.badRequest(ctx)
	}
	result, err := c.UseCase.FindByID(ctx, id.String())
	if err != nil {
		return err
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: result,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *DisbursementController) badRequest(ctx *fiber.Ctx) error {
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
	})
	ctx.Status(fiber.StatusBadRequest)
	return fiber.ErrBadRequest
}
