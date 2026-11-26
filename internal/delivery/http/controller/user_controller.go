package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/middleware"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/usecase"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UseCase *usecase.UserUseCase
	Log     *logrus.Logger
}

func NewUserController(useCase *usecase.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{UseCase: useCase, Log: log}
}

func (c *UserController) requireSuperAdmin(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	for _, role := range auth.Roles {
		if role == constants.RoleSuperAdmin {
			return nil
		}
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, constants.ErrForbiddenAccess, fiber.ErrForbidden),
	})
	ctx.Status(fiber.StatusForbidden)
	return fiber.ErrForbidden
}

func (c *UserController) ListAdmins(ctx *fiber.Ctx) error {
	if err := c.requireSuperAdmin(ctx); err != nil {
		return err
	}
	request := new(payload.QueryAdminRequest)
	if err := ctx.QueryParser(request); err != nil {
		return c.badRequest(ctx)
	}
	response, err := c.UseCase.FindAdmins(ctx, request)
	if err != nil {
		return err
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{Status: constants.ResponseNameSuccess, ResponseData: response})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *UserController) GetAdmin(ctx *fiber.Ctx) error {
	if err := c.requireSuperAdmin(ctx); err != nil {
		return err
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return c.badRequest(ctx)
	}
	response, err := c.UseCase.FindAdmin(ctx, id)
	if err != nil {
		return err
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{Status: constants.ResponseNameSuccess, ResponseData: response})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *UserController) CreateAdmin(ctx *fiber.Ctx) error {
	if err := c.requireSuperAdmin(ctx); err != nil {
		return err
	}

	request := new(payload.CreateAdminRequest)
	if err := ctx.BodyParser(request); err != nil {
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.CreateAdmin(ctx, request)
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

func (c *UserController) UpdateAdmin(ctx *fiber.Ctx) error {
	if err := c.requireSuperAdmin(ctx); err != nil {
		return err
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return c.badRequest(ctx)
	}
	request := new(payload.UpdateAdminRequest)
	if err := ctx.BodyParser(request); err != nil {
		return c.badRequest(ctx)
	}
	response, err := c.UseCase.UpdateAdmin(ctx, id, request)
	if err != nil {
		return err
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{Status: constants.ResponseNameSuccess, ResponseData: response})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *UserController) DeleteAdmin(ctx *fiber.Ctx) error {
	if err := c.requireSuperAdmin(ctx); err != nil {
		return err
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return c.badRequest(ctx)
	}
	response, err := c.UseCase.DeleteAdmin(ctx, id)
	if err != nil {
		return err
	}
	ctx.Locals(constants.Response, payload.ResponseInterface{Status: constants.ResponseNameSuccess, ResponseData: response})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *UserController) badRequest(ctx *fiber.Ctx) error {
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
	})
	ctx.Status(fiber.StatusBadRequest)
	return fiber.ErrBadRequest
}
