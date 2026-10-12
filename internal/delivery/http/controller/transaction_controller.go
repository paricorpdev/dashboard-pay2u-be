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

type TransactionController struct {
	UseCase *usecase.TransactionUseCase
	Log     *logrus.Logger
}

func NewTransactionController(useCase *usecase.TransactionUseCase, log *logrus.Logger) *TransactionController {
	return &TransactionController{
		UseCase: useCase,
		Log:     log,
	}
}

// FindAll handles paginated and filtered transaction search.
func (c *TransactionController) FindAll(ctx *fiber.Ctx) error {
	request := new(payload.QueryTransactionRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse transaction query params: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.FindAll(ctx, request)
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

// FindByID handles retrieving transaction details by UUID.
func (c *TransactionController) FindByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Invalid transaction id param: %s", ctx.Params("id"))
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.FindByID(ctx, id)
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
