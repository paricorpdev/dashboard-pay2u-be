package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
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

type MerchantUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	MerchantRepository *repository.MerchantRepository
	UserRepository     *repository.UserRepository
	MS                 *msgateway.MerchantClient
}

func NewMerchantUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	merchantRepo *repository.MerchantRepository,
	userRepo *repository.UserRepository,
	ms *msgateway.MerchantClient,
) *MerchantUseCase {
	return &MerchantUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		MerchantRepository: merchantRepo,
		UserRepository:     userRepo,
		MS:                 ms,
	}
}

func (u *MerchantUseCase) GetOptions(ctx *fiber.Ctx) ([]payload.MerchantOption, error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	active := true
	options := make([]payload.MerchantOption, 0)
	for page := 1; ; page++ {
		merchants, err := u.MS.List(ctx.UserContext(), masterKey, msgateway.MerchantListQuery{
			Page: page, PerPage: 100, IsActive: &active,
		})
		if err != nil {
			return nil, u.remoteError(ctx, "list merchant options", err)
		}
		u.syncMirrors(ctx, merchants.Data)
		for _, merchant := range merchants.Data {
			options = append(options, payload.MerchantOption{ID: merchant.ID, Name: merchant.Name})
		}
		if int64(page) >= merchants.Paging.TotalPage || len(merchants.Data) == 0 {
			break
		}
	}
	return options, nil
}

func (u *MerchantUseCase) Create(ctx *fiber.Ctx, request *payload.CreateMerchantRequest) (*payload.MerchantResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}

	merchant, err := u.MS.Create(ctx.UserContext(), masterKey, *request)
	if err != nil {
		return nil, u.remoteError(ctx, "create merchant", err)
	}
	u.syncMirror(ctx, merchant)
	return merchant, nil
}

func (u *MerchantUseCase) RotateAPIKey(ctx *fiber.Ctx, id uuid.UUID) (*payload.MerchantResponse, error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	merchant, err := u.MS.RotateAPIKey(ctx.UserContext(), masterKey, id.String())
	if err != nil {
		return nil, u.remoteError(ctx, "rotate merchant API key", err)
	}
	u.syncMirror(ctx, merchant)
	return merchant, nil
}

func (u *MerchantUseCase) RotateWebhookSecret(ctx *fiber.Ctx, id uuid.UUID) (*payload.MerchantResponse, error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	merchant, err := u.MS.RotateWebhookSecret(ctx.UserContext(), masterKey, id.String())
	if err != nil {
		return nil, u.remoteError(ctx, "rotate merchant webhook secret", err)
	}
	u.syncMirror(ctx, merchant)
	return merchant, nil
}

func (u *MerchantUseCase) FindAll(ctx *fiber.Ctx, request *payload.QueryMerchantRequest) (*payload.PageResponse[payload.MerchantResponse], error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	page := request.Page
	if page < 1 {
		page = 1
	}
	perPage := request.Limit
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	merchants, err := u.MS.List(ctx.UserContext(), masterKey, msgateway.MerchantListQuery{
		Page: page, PerPage: perPage, Search: request.Search,
		IsActive: request.IsActive, SortBy: request.SortBy, SortDir: request.SortDir,
	})
	if err != nil {
		return nil, u.remoteError(ctx, "list merchants", err)
	}
	u.syncMirrors(ctx, merchants.Data)

	return &payload.PageResponse[payload.MerchantResponse]{
		Data: merchants.Data,
		PageMetadata: payload.PageMetadata{
			Page:      merchants.Paging.Page,
			PerPage:   merchants.Paging.PerPage,
			TotalData: merchants.Paging.TotalData,
			TotalPage: merchants.Paging.TotalPage,
		},
	}, nil
}

func (u *MerchantUseCase) FindByID(ctx *fiber.Ctx, id uuid.UUID) (*payload.MerchantResponse, error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	merchant, err := u.MS.Get(ctx.UserContext(), masterKey, id.String())
	if err != nil {
		return nil, u.remoteError(ctx, "get merchant", err)
	}
	u.syncMirror(ctx, merchant)
	return merchant, nil
}

func (u *MerchantUseCase) Update(ctx *fiber.Ctx, id uuid.UUID, request *payload.UpdateMerchantRequest) (*payload.MerchantResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}

	webhookURL := ""
	if request.WebhookUrl != nil {
		webhookURL = *request.WebhookUrl
	}
	merchant, err := u.MS.Update(ctx.UserContext(), masterKey, id.String(), map[string]any{
		"name":        request.Name,
		"webhook_url": webhookURL,
		"is_active":   request.IsActive,
	})
	if err != nil {
		return nil, u.remoteError(ctx, "update merchant", err)
	}
	u.syncMirror(ctx, merchant)
	return merchant, nil
}

func (u *MerchantUseCase) Delete(ctx *fiber.Ctx, id uuid.UUID) (*payload.MerchantResponse, error) {
	masterKey, err := u.currentMSMasterKey(ctx)
	if err != nil {
		return nil, err
	}
	active := false
	merchant, err := u.MS.Update(ctx.UserContext(), masterKey, id.String(), map[string]any{"is_active": &active})
	if err != nil {
		return nil, u.remoteError(ctx, "deactivate merchant", err)
	}
	u.syncMirror(ctx, merchant)
	return merchant, nil
}

func (u *MerchantUseCase) syncMirrors(ctx *fiber.Ctx, merchants []payload.MerchantResponse) {
	for i := range merchants {
		u.syncMirror(ctx, &merchants[i])
	}
}

func (u *MerchantUseCase) syncMirror(ctx *fiber.Ctx, merchant *payload.MerchantResponse) {
	if merchant == nil {
		return
	}
	id, err := uuid.Parse(merchant.ID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Warnf("MS returned invalid merchant id %q", merchant.ID)
		return
	}

	createdAt, err := time.Parse(time.RFC3339Nano, merchant.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, merchant.UpdatedAt)
	if err != nil {
		updatedAt = createdAt
	}
	var adminID *uuid.UUID
	if merchant.AdminID != nil && *merchant.AdminID != "" {
		parsedAdminID, err := uuid.Parse(*merchant.AdminID)
		if err != nil {
			logger.WithContext(ctx, u.Log).Warnf("MS returned invalid merchant admin id %q", *merchant.AdminID)
			return
		}
		adminID = &parsedAdminID
	}
	mirror := &entity.Merchant{
		ID:            id,
		Code:          merchant.Code,
		Name:          merchant.Name,
		AdminID:       adminID,
		ApiKey:        merchant.ApiKey,
		WebhookUrl:    merchant.WebhookUrl,
		WebhookSecret: merchant.WebhookSecret,
		IsActive:      merchant.IsActive,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	if err := u.MerchantRepository.UpsertFromMS(u.DB, mirror); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Failed syncing merchant projection id=%s: %+v", merchant.ID, err)
	}
}

func (u *MerchantUseCase) remoteError(ctx *fiber.Ctx, operation string, err error) error {
	logger.WithContext(ctx, u.Log).Errorf("MS merchant service %s failed: %+v", operation, err)

	status := fiber.StatusBadGateway
	message := "merchant service unavailable"
	var apiErr *msgateway.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case fiber.StatusBadRequest, fiber.StatusUnauthorized, fiber.StatusForbidden, fiber.StatusNotFound, fiber.StatusConflict, fiber.StatusUnprocessableEntity:
			status = apiErr.StatusCode
			message = apiErr.Message
		}
	}

	fiberErr := fiber.NewError(status, message)
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, fmt.Sprintf("%s: %s", operation, message), fiberErr),
	})
	ctx.Status(status)
	return fiberErr
}

func (u *MerchantUseCase) currentMSMasterKey(ctx *fiber.Ctx) (string, error) {
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

func (u *MerchantUseCase) localError(ctx *fiber.Ctx, status int, message string) error {
	fiberErr := fiber.NewError(status, message)
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, message, fiberErr),
	})
	ctx.Status(status)
	return fiberErr
}
