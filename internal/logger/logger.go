package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func WithContext(ctx *fiber.Ctx, log *logrus.Logger) *logrus.Entry {
	if fields, ok := ctx.Locals("commonFields").(map[string]interface{}); ok && fields != nil {
		return log.WithFields(fields)
	}
	return log.WithFields(logrus.Fields{})
}
