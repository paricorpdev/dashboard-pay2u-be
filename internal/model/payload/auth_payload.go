package payload

import "github.com/google/uuid"

type Auth struct {
	UUID  string   `json:"uuid"`
	Roles []string `json:"roles"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginUserResponse struct {
	TokenType    string        `json:"token_type"`
	AccessToken  string        `json:"access_token"`
	ExpiredAt    int64         `json:"expired_at"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}

type VerifyUserRequest struct {
	AccessToken string `validate:"required,jwt"`
}

type LogoutUserRequest struct {
	UUID string `json:"id" validate:"required,max=100"`
}

type UserClaims struct {
	UUID  uuid.UUID `json:"uuid"`
	Roles []string  `json:"roles"`
}
