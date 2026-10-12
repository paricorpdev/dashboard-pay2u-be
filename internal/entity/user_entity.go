package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string     `gorm:"column:username;not null;uniqueIndex"`
	Email        *string    `gorm:"column:email;uniqueIndex"`
	PasswordHash string     `gorm:"column:password_hash;not null"`
	FullName     string     `gorm:"column:full_name;not null"`
	Status       string     `gorm:"column:status;not null;default:ACTIVE"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`

	Roles []*Role `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoleID"`
}

func (u *User) TableName() string {
	return "users"
}
