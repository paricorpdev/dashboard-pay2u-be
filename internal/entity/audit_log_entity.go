package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      *uuid.UUID `gorm:"column:user_id"`
	Action      string     `gorm:"column:action;not null"`
	Entity      string     `gorm:"column:entity"`
	EntityID    *uuid.UUID `gorm:"column:entity_id"`
	Description string     `gorm:"column:description"`
	IPAddress   string     `gorm:"column:ip_address"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
}

func (a *AuditLog) TableName() string {
	return "audit_logs"
}
