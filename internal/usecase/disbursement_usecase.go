package usecase

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/client/msgateway"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// DisbursementUseCase forwards dashboard requests to ms-gateway, which remains
// the source of truth for disbursement records.
type DisbursementUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
	MS             *msgateway.MerchantClient
}

func NewDisbursementUseCase(db *gorm.DB, log *logrus.Logger, userRepo *repository.UserRepository, ms *msgateway.MerchantClient) *DisbursementUseCase {
	return &DisbursementUseCase{DB: db, Log: log, UserRepository: userRepo, MS: ms}
}

func (u *DisbursementUseCase) FindAll(ctx *fiber.Ctx, request *payload.QueryDisbursementRequest) (*payload.PageResponse[payload.DisbursementResponse], error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	page := request.Page
	if page < 1 {
		page = 1
	}
	limit := request.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	result, err := u.MS.ListDisbursements(ctx.UserContext(), masterKey, msgateway.DisbursementListQuery{
		Page:       page,
		PerPage:    limit,
		Search:     request.Search,
		Status:     request.Status,
		MerchantID: request.MerchantID,
		Month:      request.Month,
		Year:       request.Year,
	})
	if err != nil {
		return nil, u.remoteError(ctx, "list disbursements", err)
	}
	return result, nil
}

func (u *DisbursementUseCase) FindByID(ctx *fiber.Ctx, id string) (*payload.DisbursementDetailResponse, error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	result, err := u.MS.GetDisbursement(ctx.UserContext(), masterKey, id)
	if err != nil {
		return nil, u.remoteError(ctx, "get disbursement", err)
	}
	return result, nil
}

func (u *DisbursementUseCase) currentMSMasterKey(ctx *fiber.Ctx) (string, error) {
	auth, ok := ctx.Locals(constants.AuthUserLocalKey).(payload.Auth)
	if !ok || auth.UUID == "" {
		return "", u.localError(ctx, fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}
	userID, err := uuid.Parse(auth.UUID)
	if err != nil {
		return "", u.localError(ctx, fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}
	user := new(entity.User)
	if err := u.UserRepository.FindByID(u.DB, user, userID, nil); err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed loading dashboard user MS key: %+v", err)
		return "", u.localError(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	if user.MSMasterKey == nil || *user.MSMasterKey == "" {
		return "", u.localError(ctx, fiber.StatusForbidden, "MS key is not provisioned for this user")
	}
	return *user.MSMasterKey, nil
}

func (u *DisbursementUseCase) remoteError(ctx *fiber.Ctx, operation string, err error) error {
	logger.WithContext(ctx, u.Log).Errorf("MS disbursement service %s failed: %+v", operation, err)
	status := fiber.StatusBadGateway
	message := "disbursement service unavailable"
	var apiErr *msgateway.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case fiber.StatusBadRequest, fiber.StatusUnauthorized, fiber.StatusForbidden, fiber.StatusNotFound, fiber.StatusConflict, fiber.StatusUnprocessableEntity:
			status = apiErr.StatusCode
			message = apiErr.Message
		}
	}
	return u.localError(ctx, status, fmt.Sprintf("%s: %s", operation, message))
}

func (u *DisbursementUseCase) localError(ctx *fiber.Ctx, status int, message string) error {
	fiberErr := fiber.NewError(status, message)
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, message, fiberErr),
	})
	ctx.Status(status)
	return fiberErr
}
