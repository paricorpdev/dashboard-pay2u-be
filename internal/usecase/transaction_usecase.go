package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/converter"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	TransactionRepository *repository.TransactionRepository
}

func NewTransactionUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	trxRepo *repository.TransactionRepository,
) *TransactionUseCase {
	return &TransactionUseCase{
		DB:                    db,
		Log:                   log,
		TransactionRepository: trxRepo,
	}
}

// FindAll returns a paginated and filtered list of transactions.
func (u *TransactionUseCase) FindAll(ctx *fiber.Ctx, req *payload.QueryTransactionRequest) (*payload.PageResponse[payload.TransactionResponse], error) {
	var transactions []entity.Transaction
	paging, err := u.TransactionRepository.FindAllWithFilter(u.DB, req, &transactions)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed querying transactions: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return &payload.PageResponse[payload.TransactionResponse]{
		Data:         converter.TransactionsToResponses(transactions),
		PageMetadata: *paging,
	}, nil
}

// FindByID returns detailed transaction info including merchant and webhook dispatches.
func (u *TransactionUseCase) FindByID(ctx *fiber.Ctx, id uuid.UUID) (*payload.TransactionDetailResponse, error) {
	var trx entity.Transaction
	if err := u.TransactionRepository.FindByID(u.DB, &trx, id); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Transaction not found: %s", id.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrTransactionNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	return converter.TransactionToDetailResponse(&trx), nil
}
