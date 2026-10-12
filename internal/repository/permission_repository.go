package repository

import (
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
