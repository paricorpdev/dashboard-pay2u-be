package repository

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	Repository[entity.AuditLog]
	Log *logrus.Logger
}

func NewAuditLogRepository(log *logrus.Logger) *AuditLogRepository {
	return &AuditLogRepository{
		Log: log,
	}
}

func (r *AuditLogRepository) Create(db *gorm.DB, auditLog *entity.AuditLog) error {
	return db.Create(auditLog).Error
}
