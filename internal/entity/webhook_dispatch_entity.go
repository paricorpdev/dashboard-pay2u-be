package entity

import (
	"time"

	"github.com/google/uuid"
)

type WebhookDispatch struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TransactionID uuid.UUID  `gorm:"column:transaction_id;not null"`
	MerchantID    uuid.UUID  `gorm:"column:merchant_id;not null"`
	TargetUrl     string     `gorm:"column:target_url;not null"`
	EventType     string     `gorm:"column:event_type;not null;default:'payment.success'"`
	Payload       string     `gorm:"column:payload;not null;type:jsonb"`
	Status        string     `gorm:"column:status;not null;default:'PENDING'"`
	Attempts      int        `gorm:"column:attempts;not null;default:0"`
	MaxAttempts   int        `gorm:"column:max_attempts;not null;default:3"`
	NextRetryAt   *time.Time `gorm:"column:next_retry_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (w *WebhookDispatch) TableName() string {
	return "webhook_dispatches"
}
