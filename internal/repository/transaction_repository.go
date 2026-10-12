package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
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
