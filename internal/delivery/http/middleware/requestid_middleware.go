package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	ContextKeyRequestID = "requestId"
	CommonFieldsKey     = "commonFields"
	HeaderRequestID     = "X-Request-Id"
)

func NewRequestIDMiddleware(log *logrus.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Get(HeaderRequestID)
		if requestId == "" {
			requestId = uuid.New().String()
		}

		commonFields := map[string]any{
			ContextKeyRequestID: requestId,
		}
		ctx.Locals(CommonFieldsKey, commonFields)
		ctx.Locals(ContextKeyRequestID, requestId)
		ctx.Set(HeaderRequestID, requestId)

		return ctx.Next()
	}
}
