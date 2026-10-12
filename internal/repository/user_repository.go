package repository

import (
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
