package middleware

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"

	"github.com/gofiber/fiber/v2"
)

func NewResponseMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Next()

		contentType := ctx.Response().Header.ContentType()
		if isFileContentType(contentType) {
			return nil
		}

		var response payload.ResponseInterface
		responseValue := ctx.Locals(constants.Response)
		isEmpty := false

		if resp, ok := responseValue.(payload.ResponseInterface); ok {
			response = resp
		} else {
			response = payload.ResponseInterface{
				Status:       "empty response",
				ResponseData: nil,
			}
			isEmpty = true
		}

		statusCode := ctx.Response().StatusCode()
		if statusCode == 0 {
			statusCode = fiber.StatusInternalServerError
		}

		if isEmpty && response.ResponseData == nil {
			return ctx.SendStatus(fiber.StatusNoContent)
		}

		if statusCode >= fiber.StatusOK && statusCode < fiber.StatusBadRequest {
			return ctx.JSON(payload.Response{
				StatusCode:   statusCode,
				Status:       response.Status,
				ResponseData: response.ResponseData,
			})
		} else {
			return ctx.JSON(payload.ResponseError{
				StatusCode: statusCode,
				Status:     response.Status,
				Error:      response.ResponseData,
				RequestId:  getRequestID(ctx),
			})
		}
	}
}

func getRequestID(ctx *fiber.Ctx) string {
	if v, ok := ctx.Locals(ContextKeyRequestID).(string); ok {
		return v
	}
	return ""
}

func isFileContentType(contentType []byte) bool {
	ct := string(contentType)
	return ct == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" ||
		ct == "application/octet-stream" ||
		ct == "application/pdf" ||
		ct == "image/png" ||
		ct == "image/jpeg" ||
		ct == "application/zip"
}
