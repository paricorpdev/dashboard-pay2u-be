package controller

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/delivery/http/middleware"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/helpers"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	UseCase *usecase.AuthUseCase
	Log     *logrus.Logger
}

func NewAuthController(useCase *usecase.AuthUseCase, log *logrus.Logger) *AuthController {
	return &AuthController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(payload.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse request body : %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return err
	}

	response, err := c.UseCase.Login(ctx, request)
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to login user : %+v", err)
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	request := new(payload.RefreshTokenRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to parse request body : %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return err
	}

	response, err := c.UseCase.Refresh(ctx, request)
	if err != nil {
		logger.WithContext(ctx, c.Log).Warnf("Failed to refresh token : %+v", err)
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	authorization := ctx.Get(fiber.HeaderAuthorization)
	token, _ := helpers.GetTokenFromBearer(authorization)

	request := &payload.LogoutUserRequest{
		UUID:  auth.UUID,
		Token: token,
	}

	response, err := c.UseCase.Logout(ctx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		return err
	}

	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameSuccess,
		ResponseData: response,
	})
	ctx.Status(fiber.StatusOK)
	return nil
}

func (c *AuthController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userUUID, err := uuid.Parse(auth.UUID)
	if err != nil {
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUnauthorized, fiber.ErrUnauthorized),
		})
		ctx.Status(fiber.StatusUnauthorized)
		return fiber.ErrUnauthorized
	}

	response, err := c.UseCase.GetUserByUUID(ctx, userUUID)
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
