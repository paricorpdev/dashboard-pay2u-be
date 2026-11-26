package repository

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByUsername(db *gorm.DB, user *entity.User, username string) error {
	return db.Where("username = ?", username).First(user).Error
}

func (r *UserRepository) FindByUsernameInclude(db *gorm.DB, user *entity.User, username string, includes []string) error {
	data := r.ApplyIncludes(db.Where("username = ?", username), includes)
	if err := data.First(user).Error; err != nil {
		r.Log.Warnf("Failed find user by username : %+v", err)
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(db *gorm.DB, user *entity.User, id uuid.UUID, includes []string) error {
	data := db.Where("id = ?", id)
	data = r.ApplyIncludes(data, includes)
	if err := data.First(user).Error; err != nil {
		r.Log.Warnf("Failed find user by id : %+v", err)
		return err
	}
	return nil
}

func (r *UserRepository) ExistsByUsername(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByEmail(db *gorm.DB, email string) (bool, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ValidateExistByUsername(db *gorm.DB, username string, id uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("username = ? AND id != ?", username, id).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ValidateExistByEmail(db *gorm.DB, email string, id uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) FindAdminByID(db *gorm.DB, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := db.Preload("Roles.Permissions").
		Where("users.id = ?", id).
		Where("EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON r.id = ur.role_id WHERE ur.user_id = users.id AND r.code = ?)", constants.RoleAdmin).
		Where("NOT EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON r.id = ur.role_id WHERE ur.user_id = users.id AND r.code = ?)", constants.RoleSuperAdmin).
		First(&user).Error
	return &user, err
}

func (r *UserRepository) FindAdmins(db *gorm.DB, request *payload.QueryAdminRequest) ([]entity.User, *payload.PageMetadata, error) {
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

	query := db.Model(&entity.User{}).
		Preload("Roles.Permissions").
		Where("EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON r.id = ur.role_id WHERE ur.user_id = users.id AND r.code = ?)", constants.RoleAdmin).
		Where("NOT EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON r.id = ur.role_id WHERE ur.user_id = users.id AND r.code = ?)", constants.RoleSuperAdmin)
	if search := strings.TrimSpace(request.Search); search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("(LOWER(username) LIKE ? OR LOWER(full_name) LIKE ? OR LOWER(email) LIKE ?)", pattern, pattern, pattern)
	}
	if status := strings.TrimSpace(request.Status); status != "" {
		query = query.Where("status = ?", strings.ToUpper(status))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}
	metadata := &payload.PageMetadata{Page: page, PerPage: perPage, TotalData: total}
	if total > 0 {
		metadata.TotalPage = int64(math.Ceil(float64(total) / float64(perPage)))
	}

	var users []entity.User
	if err := query.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&users).Error; err != nil {
		return nil, nil, err
	}
	return users, metadata, nil
}

func (r *UserRepository) FindAll(db *gorm.DB, request *payload.PageUserRequest, response *payload.WebResponse[[]entity.User]) error {

	data := db.Table("users").Model(&entity.User{})
	data = r.ApplyIncludes(data, request.Include)
	data = r.ApplySorting(data, request.Sort)

	if request.Filter.Username != "" {
		data = data.Where("LOWER(username) LIKE (?)", "%"+strings.ToLower(request.Filter.Username)+"%")
	}

	if request.Filter.Email != "" {
		data = data.Where("LOWER(email) LIKE (?)", "%"+strings.ToLower(request.Filter.Email)+"%")
	}

	if request.Filter.Status != "" {
		data = data.Where("status = ?", request.Filter.Status)
	}

	if len(request.Filter.RoleCodes) > 0 {
		data = data.Where(
			"EXISTS (SELECT 1 FROM user_roles ur JOIN roles ro ON ro.id = ur.role_id WHERE ur.user_id = users.id AND ro.code IN (?))",
			request.Filter.RoleCodes,
		)
	}

	page, perPage := request.Page, request.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	response.Paging = &payload.PageMetadata{}
	data.Count(&response.Paging.TotalData)
	if response.Paging.TotalData > 0 {
		response.Paging.TotalPage = int64(math.Ceil(float64(response.Paging.TotalData) / float64(perPage)))
	} else {
		response.Paging.TotalPage = 0
	}

	response.Paging.PerPage = perPage
	response.Paging.Page = page
	offset := (page - 1) * perPage
	if err := data.Offset(offset).Limit(perPage).Find(&response.Data).Error; err != nil {
		r.Log.Warnf("Failed find all users : %+v", err)
		return err
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(db *gorm.DB, user *entity.User) error {
	return db.Model(user).Update("last_login_at", time.Now()).Error
}

func (r *UserRepository) ResetPassword(db *gorm.DB, user *entity.User) error {
	return db.Model(user).Where("id = ?", user.ID).Update("password_hash", user.PasswordHash).Error
}

func (r *UserRepository) DeleteWithAssociations(db *gorm.DB, user *entity.User) error {
	return db.Select(clause.Associations).Delete(user).Error
}
