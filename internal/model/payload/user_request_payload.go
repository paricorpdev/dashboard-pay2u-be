package payload

import "github.com/google/uuid"

type PageUserRequest struct {
	Filter UserFilterRequest `json:"filter"`
	PageRequest
}

type UserFilterRequest struct {
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Status    string   `json:"status"`
	RoleCodes []string `json:"role_codes"`
}

type FindUserRequest struct {
	UUID    uuid.UUID `json:"uuid" form:"uuid" validate:"required,uuid"`
	Include []string  `json:"include" form:"include"`
}

type CreateUserRequest struct {
	Username  string   `json:"username" form:"username" validate:"required,min=3,max=100"`
	Password  string   `json:"password" form:"password" validate:"required,min=8,max=100"`
	Email     *string  `json:"email" form:"email" validate:"omitempty,email,max=150"`
	FullName  string   `json:"full_name" form:"full_name" validate:"required,max=150"`
	Status    string   `json:"status" form:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	RoleCodes []string `json:"role_codes" form:"role_codes" validate:"required,min=1,dive,max=50"`
}

type UpdateUserRequest struct {
	UUID      uuid.UUID `json:"uuid" form:"uuid" validate:"required,uuid"`
	Username  string    `json:"username" form:"username" validate:"max=100"`
	Password  string    `json:"password" form:"password" validate:"max=100"`
	Email     *string   `json:"email" form:"email" validate:"omitempty,email,max=150"`
	FullName  string    `json:"full_name" form:"full_name" validate:"max=150"`
	Status    string    `json:"status" form:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	RoleCodes []string  `json:"role_codes" form:"role_codes" validate:"omitempty,dive,max=50"`
}

type UserStatusRequest struct {
	UUID uuid.UUID `json:"uuid" form:"uuid" validate:"required,uuid"`
}

type DeleteUserRequest struct {
	UUID uuid.UUID `json:"uuid" form:"uuid" validate:"required,uuid"`
}

type ResetPasswordUserRequest struct {
	UUID uuid.UUID `json:"uuid" form:"uuid" validate:"required,uuid"`
}
