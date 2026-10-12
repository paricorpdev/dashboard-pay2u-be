package usecase

import (
	"encoding/json"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JwtUseCase struct {
	Config          *viper.Viper
	Log             *logrus.Logger
	Validate        *validator.Validate
	RedisRepository *repository.RedisRepository
}

func NewJwtUseCase(
	config *viper.Viper,
	log *logrus.Logger,
	validate *validator.Validate,
	redisRepository *repository.RedisRepository,
) *JwtUseCase {
	return &JwtUseCase{
		Config:          config,
		Log:             log,
		Validate:        validate,
		RedisRepository: redisRepository,
	}
}

func (u *JwtUseCase) GenerateToken(claims payload.UserClaims) (*payload.LoginUserResponse, error) {
	secretKey := u.Config.GetString("jwt.secret_key")
	issuer := u.Config.GetString("jwt.issuer")
	expiredMinutes := u.Config.GetInt64("jwt.time_expired")
	refreshExpiredMinutes := u.Config.GetInt64("jwt.refresh_time_expired")

	now := time.Now()
	expiredAt := now.Add(time.Minute * time.Duration(expiredMinutes))
	refreshExpiredAt := now.Add(time.Minute * time.Duration(refreshExpiredMinutes))

	// Access Token
	accessClaims := jwt.MapClaims{
		"uuid":  claims.UUID.String(),
		"roles": claims.Roles,
		"type":  constants.PrefixJwtAccessToken,
		"iss":   issuer,
		"exp":   expiredAt.Unix(),
		"iat":   now.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		u.Log.Warnf("Failed to sign access token: %+v", err)
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"uuid": claims.UUID.String(),
		"type": constants.PrefixJwtRefreshToken,
		"iss":  issuer,
		"exp":  refreshExpiredAt.Unix(),
		"iat":  now.Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		u.Log.Warnf("Failed to sign refresh token: %+v", err)
		return nil, err
	}

	return &payload.LoginUserResponse{
		TokenType:    constants.TokenTypeBearer,
		AccessToken:  signedAccess,
		ExpiredAt:    expiredAt.Unix(),
		RefreshToken: signedRefresh,
	}, nil
}

func (u *JwtUseCase) VerifyToken(tokenStr string) (*payload.UserClaims, error) {
	secretKey := u.Config.GetString("jwt.secret_key")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	uuidStr, _ := claims["uuid"].(string)
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, jwt.ErrSignatureInvalid
	}

	var roles []string
	if rolesRaw, ok := claims["roles"]; ok {
		rolesBytes, err := json.Marshal(rolesRaw)
		if err == nil {
			_ = json.Unmarshal(rolesBytes, &roles)
		}
	}

	return &payload.UserClaims{
		UUID:  parsedUUID,
		Roles: roles,
	}, nil
}

func (u *JwtUseCase) VerifyRefreshToken(tokenStr string) (*payload.UserClaims, error) {
	secretKey := u.Config.GetString("jwt.secret_key")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != constants.PrefixJwtRefreshToken {
		return nil, jwt.ErrSignatureInvalid
	}

	uuidStr, _ := claims["uuid"].(string)
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, jwt.ErrSignatureInvalid
	}

	return &payload.UserClaims{
		UUID: parsedUUID,
	}, nil
}
