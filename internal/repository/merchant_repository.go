package repository

import (
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MerchantRepository struct {
	Repository[entity.Merchant]
	Log *logrus.Logger
}

func NewMerchantRepository(log *logrus.Logger) *MerchantRepository {
	return &MerchantRepository{
		Log: log,
	}
}

func (r *MerchantRepository) FindByID(db *gorm.DB, merchant *entity.Merchant, id uuid.UUID) error {
	return db.Where("id = ?", id).First(merchant).Error
}

func (r *MerchantRepository) FindByCode(db *gorm.DB, merchant *entity.Merchant, code string) error {
	return db.Where("code = ?", code).First(merchant).Error
}

func (r *MerchantRepository) ExistsByCode(db *gorm.DB, code string) (bool, error) {
	var count int64
	err := db.Model(&entity.Merchant{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *MerchantRepository) ExistsByApiKey(db *gorm.DB, apiKey string) (bool, error) {
	var count int64
	err := db.Model(&entity.Merchant{}).Where("api_key = ?", apiKey).Count(&count).Error
	return count > 0, err
}

func (r *MerchantRepository) HasTransactions(db *gorm.DB, merchantID uuid.UUID) (bool, error) {
	var count int64
	err := db.Table("transactions").Where("merchant_id = ?", merchantID).Count(&count).Error
	return count > 0, err
}

func (r *MerchantRepository) FindOptions(db *gorm.DB, merchants *[]entity.Merchant) error {
	return db.Select("id, name").Where("is_active = ?", true).Order("name asc").Find(merchants).Error
}

func (r *MerchantRepository) FindAll(db *gorm.DB, request *payload.QueryMerchantRequest, merchants *[]entity.Merchant) (*payload.PageMetadata, error) {
	query := db.Model(&entity.Merchant{})

	if request.Search != "" {
		searchPattern := "%" + strings.ToLower(request.Search) + "%"
		query = query.Where("(LOWER(name) LIKE ? OR LOWER(code) LIKE ?)", searchPattern, searchPattern)
	}

	if request.IsActive != nil {
		query = query.Where("is_active = ?", *request.IsActive)
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

	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		r.Log.Warnf("Failed to count merchants: %+v", err)
		return nil, err
	}

	var totalPage int64 = 0
	if totalData > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(limit)))
	}

	sortBy := "created_at"
	if request.SortBy != "" {
		sortBy = request.SortBy
	}
	sortDir := "desc"
	if strings.ToLower(request.SortDir) == "asc" {
		sortDir = "asc"
	}
	query = query.Order(sortBy + " " + sortDir)

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(merchants).Error; err != nil {
		r.Log.Warnf("Failed to query merchants: %+v", err)
		return nil, err
	}

	return &payload.PageMetadata{
		Page:      page,
		PerPage:   limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}, nil
}
