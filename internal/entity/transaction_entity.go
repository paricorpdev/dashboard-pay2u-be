package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequestID        string     `gorm:"column:request_id;not null;uniqueIndex"`
	Provider         string     `gorm:"column:provider;not null;default:'pay2u'"`
	ProviderToken    *string    `gorm:"column:provider_token"`
	MerchantReff     string     `gorm:"column:merchant_reff;not null"`
	PaymentMethod    string     `gorm:"column:payment_method;not null"`
	PaymentCode      *string    `gorm:"column:payment_code"`
	Status           string     `gorm:"column:status;not null;default:'PENDING'"`
	Amount           int64      `gorm:"column:amount;not null"`
	AmountAdmin      int64      `gorm:"column:amount_admin;not null;default:0"`
	AmountDiscount   int64      `gorm:"column:amount_discount;not null;default:0"`
	AmountTotal      int64      `gorm:"column:amount_total;not null"`
	Currency         string     `gorm:"column:currency;not null;default:'IDR'"`
	CustomerName     *string    `gorm:"column:customer_name"`
	CustomerPhone    *string    `gorm:"column:customer_phone"`
	CustomerEmail    *string    `gorm:"column:customer_email"`
	BillTitle        *string    `gorm:"column:bill_title"`
	BillDescription  *string    `gorm:"column:bill_description"`
	PaymentReff      *string    `gorm:"column:payment_reff"`
	ExpiredAt        *time.Time `gorm:"column:expired_at"`
	PaidAt           *time.Time `gorm:"column:paid_at"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	MerchantID       *uuid.UUID `gorm:"column:merchant_id"`
	ProviderCallback *string    `gorm:"column:provider_callback"`

	Merchant          *Merchant         `gorm:"foreignKey:MerchantID;references:ID"`
	WebhookDispatches []WebhookDispatch `gorm:"foreignKey:TransactionID;references:ID"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}
