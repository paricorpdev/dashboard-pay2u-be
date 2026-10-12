package repository

import (
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Log: log,
	}
}

type MonthlySummaryResult struct {
	TotalCount    int64
	TotalVolume   int64
	SuccessCount  int64
	PendingCount  int64
	ExpiredCount  int64
	TotalAdminFee int64
}

type RawDailyTrend struct {
	DayStr       string `gorm:"column:day_str"`
	TotalCount   int64  `gorm:"column:total_count"`
	TotalAmount  int64  `gorm:"column:total_amount"`
	SuccessCount int64  `gorm:"column:success_count"`
	PendingCount int64  `gorm:"column:pending_count"`
	ExpiredCount int64  `gorm:"column:expired_count"`
}

type RawPaymentMethodDist struct {
	PaymentMethod string `gorm:"column:payment_method"`
	Count         int64  `gorm:"column:count"`
	Amount        int64  `gorm:"column:amount"`
}

type RawStatusDist struct {
	Status string `gorm:"column:status"`
	Count  int64  `gorm:"column:count"`
	Amount int64  `gorm:"column:amount"`
}

func (r *TransactionRepository) FindByID(db *gorm.DB, trx *entity.Transaction, id uuid.UUID) error {
	return db.Preload("Merchant").Preload("WebhookDispatches").Where("id = ?", id).First(trx).Error
}

func (r *TransactionRepository) FindAllWithFilter(db *gorm.DB, req *payload.QueryTransactionRequest, transactions *[]entity.Transaction) (*payload.PageMetadata, error) {
	query := db.Model(&entity.Transaction{})

	if req.Name != "" {
		query = query.Where("LOWER(customer_name) LIKE ?", "%"+strings.ToLower(req.Name)+"%")
	}
	if req.Phone != "" {
		query = query.Where("customer_phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Email != "" {
		query = query.Where("LOWER(customer_email) LIKE ?", "%"+strings.ToLower(req.Email)+"%")
	}
	if req.PaymentMethod != "" {
		query = query.Where("payment_method = ?", req.PaymentMethod)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.MerchantReff != "" {
		query = query.Where("LOWER(merchant_reff) LIKE ?", "%"+strings.ToLower(req.MerchantReff)+"%")
	}
	if req.BillTitle != "" {
		query = query.Where("LOWER(bill_title) LIKE ?", "%"+strings.ToLower(req.BillTitle)+"%")
	}
	if req.StatusPaid != nil {
		if *req.StatusPaid {
			query = query.Where("paid_at IS NOT NULL")
		} else {
			query = query.Where("paid_at IS NULL")
		}
	}
	if req.MerchantID != nil && *req.MerchantID != "" {
		if mID, err := uuid.Parse(*req.MerchantID); err == nil {
			query = query.Where("merchant_id = ?", mID)
		}
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		r.Log.Warnf("Failed to count transactions: %+v", err)
		return nil, err
	}

	var totalPage int64 = 0
	if totalData > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(limit)))
	}

	sortBy := "created_at"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	sortDir := "desc"
	if strings.ToLower(req.SortDir) == "asc" {
		sortDir = "asc"
	}
	query = query.Order(sortBy + " " + sortDir)

	offset := (page - 1) * limit
	if err := query.Preload("Merchant").Offset(offset).Limit(limit).Find(transactions).Error; err != nil {
		r.Log.Warnf("Failed to find transactions: %+v", err)
		return nil, err
	}

	return &payload.PageMetadata{
		Page:      page,
		PerPage:   limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}, nil
}

func (r *TransactionRepository) GetMonthlySummary(db *gorm.DB, start time.Time, end time.Time, merchantID *uuid.UUID) (*MonthlySummaryResult, error) {
	var res MonthlySummaryResult
	query := db.Table("transactions").
		Select(`
			COUNT(*) AS total_count,
			COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN amount_total ELSE 0 END), 0) AS total_volume,
			COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END), 0) AS success_count,
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending_count,
			COALESCE(SUM(CASE WHEN status = 'EXPIRED' THEN 1 ELSE 0 END), 0) AS expired_count,
			COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN amount_admin ELSE 0 END), 0) AS total_admin_fee
		`).
		Where("created_at >= ? AND created_at <= ?", start, end)

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}

	if err := query.Scan(&res).Error; err != nil {
		r.Log.Warnf("Failed to query monthly summary: %+v", err)
		return nil, err
	}
	return &res, nil
}

func (r *TransactionRepository) GetDailyTrend(db *gorm.DB, start time.Time, end time.Time, merchantID *uuid.UUID) ([]RawDailyTrend, error) {
	var rows []RawDailyTrend
	query := db.Table("transactions").
		Select(`
			TO_CHAR(created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') AS day_str,
			COUNT(*) AS total_count,
			COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN amount_total ELSE 0 END), 0) AS total_amount,
			COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END), 0) AS success_count,
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending_count,
			COALESCE(SUM(CASE WHEN status = 'EXPIRED' THEN 1 ELSE 0 END), 0) AS expired_count
		`).
		Where("created_at >= ? AND created_at <= ?", start, end)

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}

	query = query.Group("day_str").Order("day_str ASC")

	if err := query.Scan(&rows).Error; err != nil {
		r.Log.Warnf("Failed to query daily trend: %+v", err)
		return nil, err
	}
	return rows, nil
}

func (r *TransactionRepository) GetPaymentMethodDistribution(db *gorm.DB, start time.Time, end time.Time, merchantID *uuid.UUID) ([]RawPaymentMethodDist, error) {
	var rows []RawPaymentMethodDist
	query := db.Table("transactions").
		Select(`
			payment_method,
			COUNT(*) AS count,
			COALESCE(SUM(amount_total), 0) AS amount
		`).
		Where("created_at >= ? AND created_at <= ?", start, end)

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}

	query = query.Group("payment_method").Order("count DESC")

	if err := query.Scan(&rows).Error; err != nil {
		r.Log.Warnf("Failed to query payment method distribution: %+v", err)
		return nil, err
	}
	return rows, nil
}

func (r *TransactionRepository) GetStatusDistribution(db *gorm.DB, start time.Time, end time.Time, merchantID *uuid.UUID) ([]RawStatusDist, error) {
	var rows []RawStatusDist
	query := db.Table("transactions").
		Select(`
			status,
			COUNT(*) AS count,
			COALESCE(SUM(amount_total), 0) AS amount
		`).
		Where("created_at >= ? AND created_at <= ?", start, end)

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}

	query = query.Group("status").Order("count DESC")

	if err := query.Scan(&rows).Error; err != nil {
		r.Log.Warnf("Failed to query status distribution: %+v", err)
		return nil, err
	}
	return rows, nil
}
