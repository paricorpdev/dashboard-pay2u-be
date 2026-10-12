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

type RBACController struct {
	UseCase *usecase.RBACUseCase
	Log     *logrus.Logger
}

func NewRBACController(useCase *usecase.RBACUseCase, log *logrus.Logger) *RBACController {
	return &RBACController{
		UseCase: useCase,
		Log:     log,
	}
}

// FindAllRoles returns all roles in the system.
func (c *RBACController) FindAllRoles(ctx *fiber.Ctx) error {
	response, err := c.UseCase.FindAllRoles(ctx)
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

// FindRoleByID returns details and permissions of a selected role.
func (c *RBACController) FindRoleByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Invalid role id param: %s", ctx.Params("id"))
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInvalidRequest, fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.FindRoleByID(ctx, id)
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

// FindAllPermissions returns all system permissions (flat or grouped by resource).
func (c *RBACController) FindAllPermissions(ctx *fiber.Ctx) error {
	isGrouped := ctx.Query("grouped") == "true"

	var response any
	var err error

	if isGrouped {
		response, err = c.UseCase.FindPermissionsGrouped(ctx)
	} else {
		response, err = c.UseCase.FindAllPermissions(ctx)
	}

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
