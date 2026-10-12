package usecase

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/helpers"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/converter"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AnalyticUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	TransactionRepository *repository.TransactionRepository
	TimeLocation          *time.Location
}

func NewAnalyticUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	trxRepo *repository.TransactionRepository,
	timeZone string,
) *AnalyticUseCase {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		loc = time.Local
	}
	return &AnalyticUseCase{
		DB:                    db,
		Log:                   log,
		TransactionRepository: trxRepo,
		TimeLocation:          loc,
	}
}

func (u *AnalyticUseCase) parseMerchantID(merchantIDStr *string) *uuid.UUID {
	if merchantIDStr == nil || *merchantIDStr == "" {
		return nil
	}
	parsed, err := uuid.Parse(*merchantIDStr)
	if err != nil {
		return nil
	}
	return &parsed
}

func (u *AnalyticUseCase) GetSummary(ctx *fiber.Ctx, request *payload.AnalyticFilterRequest) (*payload.AnalyticSummaryResponse, error) {
	month, year := helpers.ResolveMonthYear(request.Month, request.Year)
	startDate, endDate, _ := helpers.GetMonthDateRange(year, month, u.TimeLocation)
	merchantUUID := u.parseMerchantID(request.MerchantID)

	summaryResult, err := u.TransactionRepository.GetMonthlySummary(u.DB, startDate, endDate, merchantUUID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to query monthly summary: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SummaryToResponse(month, year, summaryResult), nil
}

func (u *AnalyticUseCase) GetDailyTrend(ctx *fiber.Ctx, request *payload.AnalyticFilterRequest) ([]payload.DailyTrendItem, error) {
	month, year := helpers.ResolveMonthYear(request.Month, request.Year)
	startDate, endDate, daysInMonth := helpers.GetMonthDateRange(year, month, u.TimeLocation)
	merchantUUID := u.parseMerchantID(request.MerchantID)

	rawTrends, err := u.TransactionRepository.GetDailyTrend(u.DB, startDate, endDate, merchantUUID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to query daily trend: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RawDailyTrendsToCompleteMonth(year, month, daysInMonth, rawTrends), nil
}

func (u *AnalyticUseCase) GetPaymentMethodDistribution(ctx *fiber.Ctx, request *payload.AnalyticFilterRequest) ([]payload.PaymentMethodDistributionItem, error) {
	month, year := helpers.ResolveMonthYear(request.Month, request.Year)
	startDate, endDate, _ := helpers.GetMonthDateRange(year, month, u.TimeLocation)
	merchantUUID := u.parseMerchantID(request.MerchantID)

	summaryResult, err := u.TransactionRepository.GetMonthlySummary(u.DB, startDate, endDate, merchantUUID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to query summary for payment distribution: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	rawDistribution, err := u.TransactionRepository.GetPaymentMethodDistribution(u.DB, startDate, endDate, merchantUUID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to query payment method distribution: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RawPaymentMethodToDistribution(rawDistribution, summaryResult.TotalCount), nil
}

func (u *AnalyticUseCase) GetStatusDistribution(ctx *fiber.Ctx, request *payload.AnalyticFilterRequest) ([]payload.StatusDistributionItem, error) {
	month, year := helpers.ResolveMonthYear(request.Month, request.Year)
	startDate, endDate, _ := helpers.GetMonthDateRange(year, month, u.TimeLocation)
	merchantUUID := u.parseMerchantID(request.MerchantID)

	summaryResult, err := u.TransactionRepository.GetMonthlySummary(u.DB, startDate, endDate, merchantUUID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to query summary for status distribution: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	rawDistribution, err := u.TransactionRepository.GetStatusDistribution(u.DB, startDate, endDate, merchantUUID)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to query status distribution: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RawStatusToDistribution(rawDistribution, summaryResult.TotalCount), nil
}
