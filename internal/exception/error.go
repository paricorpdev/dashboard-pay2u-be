package exception

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHandleError(ctx *fiber.Ctx, details interface{}, fiberError *fiber.Error) fiber.Map {
	errorResponse := fiber.Map{
		"message":   fiberError.Message,
		"details":   details,
		"timestamp": time.Now().Format(time.RFC3339),
		"path":      ctx.Path(),
	}
	return errorResponse
}
