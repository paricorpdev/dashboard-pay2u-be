package entity

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code          string    `gorm:"column:code;not null;uniqueIndex"`
	Name          string    `gorm:"column:name;not null"`
	ApiKey        string    `gorm:"column:api_key;not null;uniqueIndex"`
	WebhookUrl    *string   `gorm:"column:webhook_url"`
	WebhookSecret *string   `gorm:"column:webhook_secret"`
	IsActive      bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (m *Merchant) TableName() string {
	return "merchants"
}
