package repository

import (
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleRepository struct {
	Repository[entity.Role]
	Log *logrus.Logger
}

func NewRoleRepository(log *logrus.Logger) *RoleRepository {
	return &RoleRepository{
		Log: log,
	}
}

// FindAll returns all roles for lightweight listing.
func (r *RoleRepository) FindAll(db *gorm.DB) ([]entity.Role, error) {
	var roles []entity.Role
	err := db.Order("code ASC").Find(&roles).Error
	if err != nil {
		r.Log.Warnf("Failed to query all roles: %+v", err)
		return nil, err
	}
	return roles, nil
}

// FindByID returns a role by UUID preloading its permissions.
func (r *RoleRepository) FindByID(db *gorm.DB, role *entity.Role, id uuid.UUID) error {
	return db.Preload("Permissions").Where("id = ?", id).First(role).Error
}

// CountUsersByRoleID returns the number of active users assigned to a role.
func (r *RoleRepository) CountUsersByRoleID(db *gorm.DB, roleID uuid.UUID) (int64, error) {
	var count int64
	err := db.Table("user_roles").Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

// FindByCodes resolves role entities by their code strings.
func (r *RoleRepository) FindByCodes(db *gorm.DB, codes []string) ([]entity.Role, error) {
	var roles []entity.Role
	if len(codes) == 0 {
		return roles, nil
	}

	err := db.Where("code IN (?)", codes).Find(&roles).Error
	if err != nil {
		r.Log.Warnf("Failed find roles by codes : %+v", err)
		return nil, err
	}

	return roles, nil
}
