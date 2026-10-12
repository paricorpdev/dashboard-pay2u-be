package usecase

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/helpers"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/converter"
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
}

func NewMerchantUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	merchantRepo *repository.MerchantRepository,
) *MerchantUseCase {
	return &MerchantUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		MerchantRepository: merchantRepo,
	}
}

func (u *MerchantUseCase) Create(ctx *fiber.Ctx, request *payload.CreateMerchantRequest) (*payload.MerchantResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Invalid create merchant request: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}

	exists, err := u.MerchantRepository.ExistsByCode(u.DB, request.Code)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed checking merchant code existence: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}
	if exists {
		logger.WithContext(ctx, u.Log).Warnf("Merchant code already exists: %s", request.Code)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrMerchantCodeAlreadyExists, fiber.ErrConflict),
		})
		ctx.Status(fiber.StatusConflict)
		return nil, fiber.ErrConflict
	}

	apiKey, err := helpers.GenerateApiKey()
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed generating api key: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrFailedGenerateCredentials, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	webhookSecret, err := helpers.GenerateWebhookSecret()
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed generating webhook secret: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrFailedGenerateCredentials, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	now := time.Now()
	merchant := entity.Merchant{
		ID:            uuid.New(),
		Code:          request.Code,
		Name:          request.Name,
		ApiKey:        apiKey,
		WebhookUrl:    request.WebhookUrl,
		WebhookSecret: &webhookSecret,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	tx := u.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	if err := u.MerchantRepository.Create(tx, &merchant); err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed creating merchant: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	tx.Commit()

	return converter.MerchantToResponse(&merchant), nil
}

func (u *MerchantUseCase) FindAll(ctx *fiber.Ctx, request *payload.QueryMerchantRequest) (*payload.PageResponse[payload.MerchantResponse], error) {
	var merchants []entity.Merchant
	paging, err := u.MerchantRepository.FindAll(u.DB, request, &merchants)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed querying merchants: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return &payload.PageResponse[payload.MerchantResponse]{
		Data:         converter.MerchantsToResponses(merchants),
		PageMetadata: *paging,
	}, nil
}

func (u *MerchantUseCase) FindByID(ctx *fiber.Ctx, id uuid.UUID) (*payload.MerchantResponse, error) {
	var merchant entity.Merchant
	if err := u.MerchantRepository.FindByID(u.DB, &merchant, id); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Merchant not found: %s", id.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrMerchantNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	return converter.MerchantToResponse(&merchant), nil
}

func (u *MerchantUseCase) Update(ctx *fiber.Ctx, id uuid.UUID, request *payload.UpdateMerchantRequest) (*payload.MerchantResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Invalid update merchant request: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}

	tx := u.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	var merchant entity.Merchant
	if err := u.MerchantRepository.FindByID(tx, &merchant, id); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Merchant not found for update: %s", id.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrMerchantNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	merchant.Name = request.Name
	merchant.WebhookUrl = request.WebhookUrl
	if request.IsActive != nil {
		merchant.IsActive = *request.IsActive
	}
	merchant.UpdatedAt = time.Now()

	if err := u.MerchantRepository.Update(tx, &merchant); err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed updating merchant: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	tx.Commit()

	return converter.MerchantToResponse(&merchant), nil
}

func (u *MerchantUseCase) Delete(ctx *fiber.Ctx, id uuid.UUID) (*payload.MerchantResponse, error) {
	tx := u.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	var merchant entity.Merchant
	if err := u.MerchantRepository.FindByID(tx, &merchant, id); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Merchant not found for delete: %s", id.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrMerchantNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	// Soft-deactivate to preserve foreign key associations and financial transaction history
	merchant.IsActive = false
	merchant.UpdatedAt = time.Now()

	if err := u.MerchantRepository.Update(tx, &merchant); err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed deactivating merchant: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	tx.Commit()

	return converter.MerchantToResponse(&merchant), nil
}

func (u *MerchantUseCase) RegenerateSecret(ctx *fiber.Ctx, id uuid.UUID, request *payload.RegenerateMerchantSecretRequest) (*payload.MerchantResponse, error) {
	tx := u.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	var merchant entity.Merchant
	if err := u.MerchantRepository.FindByID(tx, &merchant, id); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Merchant not found for secret regeneration: %s", id.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrMerchantNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	if request.RotateApiKey {
		newKey, err := helpers.GenerateApiKey()
		if err != nil {
			logger.WithContext(ctx, u.Log).Errorf("Failed generating new api key: %+v", err)
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrFailedGenerateCredentials, fiber.ErrInternalServerError),
			})
			ctx.Status(fiber.StatusInternalServerError)
			return nil, fiber.ErrInternalServerError
		}
		merchant.ApiKey = newKey
	}

	if request.RotateWebhookSecret {
		newSecret, err := helpers.GenerateWebhookSecret()
		if err != nil {
			logger.WithContext(ctx, u.Log).Errorf("Failed generating new webhook secret: %+v", err)
			ctx.Locals(constants.Response, payload.ResponseInterface{
				Status:       constants.ResponseNameError,
				ResponseData: exception.NewHandleError(ctx, constants.ErrFailedGenerateCredentials, fiber.ErrInternalServerError),
			})
			ctx.Status(fiber.StatusInternalServerError)
			return nil, fiber.ErrInternalServerError
		}
		merchant.WebhookSecret = &newSecret
	}

	merchant.UpdatedAt = time.Now()

	if err := u.MerchantRepository.Update(tx, &merchant); err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed saving rotated merchant credentials: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	tx.Commit()

	return converter.MerchantToResponse(&merchant), nil
}
