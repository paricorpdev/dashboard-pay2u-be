package middleware

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/helpers"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const userLocalKey = constants.AuthUserLocalKey

func NewAuth(log *logrus.Logger, config *viper.Viper, jwtUseCase *usecase.JwtUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorization := ctx.Get(fiber.HeaderAuthorization)
		if authorization == "" {
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrUnauthorized, fiber.ErrUnauthorized),
			})
			ctx.Status(fiber.StatusUnauthorized)
			return fiber.ErrUnauthorized
		}

		token, err := helpers.GetTokenFromBearer(authorization)
		if err != nil {
			log.Warnf("Failed to parse bearer token : %+v", err)
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrUnauthorized, fiber.ErrUnauthorized),
			})
			ctx.Status(fiber.StatusUnauthorized)
			return fiber.ErrUnauthorized
		}

		claims, err := jwtUseCase.VerifyToken(token)
		if err != nil {
			log.Warnf("Failed to verify token : %+v", err)
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrUnauthorized, fiber.ErrUnauthorized),
			})
			ctx.Status(fiber.StatusUnauthorized)
			return fiber.ErrUnauthorized
		}

		ctx.Locals(userLocalKey, payload.Auth{
			UUID:  claims.UUID.String(),
			Roles: claims.Roles,
		})

		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) payload.Auth {
	value := ctx.Locals(userLocalKey)
	if value == nil {
		return payload.Auth{}
	}

	auth, ok := value.(payload.Auth)
	if !ok {
		return payload.Auth{}
	}

	return auth
}
