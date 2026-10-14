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

type MerchantController struct {
	UseCase *usecase.MerchantUseCase
	Log     *logrus.Logger
}

func NewMerchantController(useCase *usecase.MerchantUseCase, log *logrus.Logger) *MerchantController {
	return &MerchantController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *MerchantController) GetOptions(ctx *fiber.Ctx) error {
	options, err := c.UseCase.GetOptions(ctx)
	if err != nil {
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: options,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *MerchantController) Create(ctx *fiber.Ctx) error {
	request := new(payload.CreateMerchantRequest)
	if err := ctx.BodyParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse create merchant body: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx, request)
	if err != nil {
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusCreated)
	return nil
}

func (c *MerchantController) FindAll(ctx *fiber.Ctx) error {
	request := new(payload.QueryMerchantRequest)
	if err := ctx.QueryParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse query params: %+v", err)
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

func (c *MerchantController) FindByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Invalid merchant id param: %s", ctx.Params("id"))
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

func (c *MerchantController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Invalid merchant id param for update: %s", ctx.Params("id"))
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	request := new(payload.UpdateMerchantRequest)
	if err := ctx.BodyParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse update merchant body: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Update(ctx, id, request)
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

func (c *MerchantController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Invalid merchant id param for delete: %s", ctx.Params("id"))
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Delete(ctx, id)
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

func (c *MerchantController) RegenerateSecret(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Invalid merchant id param for secret rotation: %s", ctx.Params("id"))
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	request := new(payload.RegenerateMerchantSecretRequest)
	if err := ctx.BodyParser(request); err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse regenerate secret body: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.RegenerateSecret(ctx, id, request)
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
