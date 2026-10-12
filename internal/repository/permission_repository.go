package repository

import (
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	Repository[entity.Permission]
	Log *logrus.Logger
}

func NewPermissionRepository(log *logrus.Logger) *PermissionRepository {
	return &PermissionRepository{
		Log: log,
	}
}

// UserHasAccess checks DB whether the user has permission for a specific resource and action.
func (r *PermissionRepository) UserHasAccess(db *gorm.DB, userID string, resource string, action string) (bool, error) {
	var count int64
	err := db.Table("permissions p").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ?", userID).
		Where("p.resource = ? AND p.action = ?", resource, action).
		Count(&count).Error
	return count > 0, err
}

// FindAll returns all registered system permissions ordered by resource and action.
func (r *PermissionRepository) FindAll(db *gorm.DB) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := db.Order("resource ASC, action ASC").Find(&perms).Error
	if err != nil {
		r.Log.Warnf("Failed to query permissions: %+v", err)
		return nil, err
	}
	return perms, nil
}

// GetUserPermissions returns all distinct permissions granted to a user via their assigned roles.
func (r *PermissionRepository) GetUserPermissions(db *gorm.DB, userID uuid.UUID) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := db.Table("permissions p").
		Distinct("p.*").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ?", userID).
		Order("p.resource ASC, p.action ASC").
		Scan(&perms).Error
	if err != nil {
		r.Log.Warnf("Failed to query user permissions: %+v", err)
		return nil, err
	}
	return perms, nil
}
