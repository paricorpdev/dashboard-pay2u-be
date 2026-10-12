package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code        string    `gorm:"column:code;not null;uniqueIndex"`
	Resource    string    `gorm:"column:resource;not null"`
	Action      string    `gorm:"column:action;not null"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (p *Permission) TableName() string {
	return "permissions"
}
