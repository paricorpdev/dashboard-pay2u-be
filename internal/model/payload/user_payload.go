package payload

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID      `json:"id"`
	Username    string         `json:"username"`
	Email       *string        `json:"email"`
	FullName    string         `json:"full_name"`
	Status      string         `json:"status"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	Roles       []RoleResponse `json:"roles"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleListItemResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleDetailResponse struct {
	ID          uuid.UUID            `json:"id"`
	Code        string               `json:"code"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	UserCount   int64                `json:"user_count"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type PermissionGroupResponse struct {
	Resource    string               `json:"resource"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
