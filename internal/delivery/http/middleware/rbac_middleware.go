package middleware

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewRBACMiddleware(db *gorm.DB, log *logrus.Logger, permissionRepo *repository.PermissionRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetUser(ctx)
		if auth.UUID == "" {
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrUnauthorized, fiber.ErrUnauthorized),
			})
			ctx.Status(fiber.StatusUnauthorized)
			return fiber.ErrUnauthorized
		}

		resource := ""
		segments := strings.Split(strings.Trim(ctx.Path(), "/"), "/")
		for i, segment := range segments {
			if segment == "api" && i+2 < len(segments) && strings.HasPrefix(segments[i+1], "v") {
				resource = segments[i+2]
				break
			}
		}

		if resource == "" || resource == "auth" {
			return ctx.Next()
		}

		action := "read"
		if len(segments) > 0 {
			for _, segment := range segments {
				switch segment {
				case "create", "upload", "sync":
					action = "create"
				case "update", "reset-password":
					action = "update"
				case "activate":
					action = "activate"
				case "deactivate":
					action = "deactivate"
				case "delete":
					action = "delete"
				}
			}
		}

		hasPermission, err := permissionRepo.UserHasAccess(db, auth.UUID, resource, action)
		if err != nil {
			log.Warnf("RBAC check DB error for user=%s resource=%s action=%s: %+v", auth.UUID, resource, action, err)
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
			})
			ctx.Status(fiber.StatusInternalServerError)
			return fiber.ErrInternalServerError
		}

		if !hasPermission {
			log.Warnf("RBAC access denied. user=%s resource=%s action=%s", auth.UUID, resource, action)
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrForbiddenAccess, fiber.ErrForbidden),
			})
			ctx.Status(fiber.StatusForbidden)
			return fiber.ErrForbidden
		}

		return ctx.Next()
	}
}
