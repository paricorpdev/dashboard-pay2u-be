package repository

import (
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
