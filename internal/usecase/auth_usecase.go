package usecase

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/hashing"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/converter"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB              *gorm.DB
	Config          *viper.Viper
	Log             *logrus.Logger
	Validate        *validator.Validate
	JwtUseCase      *JwtUseCase
	UserRepository  *repository.UserRepository
	RedisRepository *repository.RedisRepository
}

func NewAuthUseCase(
	db *gorm.DB,
	config *viper.Viper,
	log *logrus.Logger,
	validate *validator.Validate,
	jwtUseCase *JwtUseCase,
	userRepository *repository.UserRepository,
	redisRepository *repository.RedisRepository,
) *AuthUseCase {
	return &AuthUseCase{
		DB:              db,
		Config:          config,
		Log:             log,
		Validate:        validate,
		JwtUseCase:      jwtUseCase,
		UserRepository:  userRepository,
		RedisRepository: redisRepository,
	}
}

func (u *AuthUseCase) Login(ctx *fiber.Ctx, request *payload.LoginUserRequest) (*payload.LoginUserResponse, error) {
	tx := u.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Invalid login request : %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByUsernameInclude(tx, user, request.Username, []string{"Roles", "Roles.Permissions"}); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("User not found: %s", request.Username)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUsernameOrPassword, fiber.ErrUnauthorized),
		})
		ctx.Status(fiber.StatusUnauthorized)
		return nil, fiber.ErrUnauthorized
	}

	if user.Status != constants.StatusActive {
		logger.WithContext(ctx, u.Log).Warnf("Inactive user tried to login: %s", request.Username)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUserInactive, fiber.ErrForbidden),
		})
		ctx.Status(fiber.StatusForbidden)
		return nil, fiber.ErrForbidden
	}

	if err := hashing.ComparePasswords(user.PasswordHash, request.Password); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Wrong password for user: %s", request.Username)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUsernameOrPassword, fiber.ErrUnauthorized),
		})
		ctx.Status(fiber.StatusUnauthorized)
		return nil, fiber.ErrUnauthorized
	}

	var roleCodes []string
	for _, r := range user.Roles {
		roleCodes = append(roleCodes, r.Code)
	}

	tokenResponse, err := u.JwtUseCase.GenerateToken(payload.UserClaims{
		UUID:  user.ID,
		Roles: roleCodes,
	})
	if err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Failed to generate token: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrFailedGenerate, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := u.UserRepository.UpdateLastLogin(tx, user); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Failed to update last login: %+v", err)
	}

	tx.Commit()

	tokenResponse.User = converter.UserToResponse(user)
	return tokenResponse, nil
}

func (u *AuthUseCase) Refresh(ctx *fiber.Ctx, request *payload.RefreshTokenRequest) (*payload.LoginUserResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Invalid refresh request : %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}

	claims, err := u.JwtUseCase.VerifyRefreshToken(request.RefreshToken)
	if err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Invalid refresh token: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUnauthorized, fiber.ErrUnauthorized),
		})
		ctx.Status(fiber.StatusUnauthorized)
		return nil, fiber.ErrUnauthorized
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(u.DB, user, claims.UUID, []string{"Roles"}); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("User not found for refresh: %s", claims.UUID.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUserNotFound, fiber.ErrUnauthorized),
		})
		ctx.Status(fiber.StatusUnauthorized)
		return nil, fiber.ErrUnauthorized
	}

	var roleCodes []string
	for _, r := range user.Roles {
		roleCodes = append(roleCodes, r.Code)
	}

	tokenResponse, err := u.JwtUseCase.GenerateToken(payload.UserClaims{
		UUID:  user.ID,
		Roles: roleCodes,
	})
	if err != nil {
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrFailedGenerate, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return tokenResponse, nil
}

func (u *AuthUseCase) Logout(ctx *fiber.Ctx, request *payload.LogoutUserRequest) (fiber.Map, error) {
	if err := u.Validate.Struct(request); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Invalid logout request : %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, err.Error(), fiber.ErrBadRequest),
		})
		ctx.Status(fiber.StatusBadRequest)
		return nil, fiber.ErrBadRequest
	}

	if u.RedisRepository != nil && u.RedisRepository.Redis != nil {
		_ = u.RedisRepository.Del(constants.PrefixJwtAccessToken + ":" + request.UUID)
		_ = u.RedisRepository.Del(constants.PrefixJwtRefreshToken + ":" + request.UUID)

		if request.Token != "" {
			expiredMinutes := u.Config.GetInt64("jwt.time_expired")
			if expiredMinutes <= 0 {
				expiredMinutes = 360
			}
			ttl := time.Minute * time.Duration(expiredMinutes)
			_ = u.RedisRepository.Set("jwt_blacklist:"+request.Token, []byte("1"), ttl)
		}
	}

	return fiber.Map{
		"message": "logged out",
	}, nil
}

// GetUserByUUID fetches the current user by UUID (for /users/current).
func (u *AuthUseCase) GetUserByUUID(ctx *fiber.Ctx, userUUID uuid.UUID) (*payload.UserResponse, error) {
	user := new(entity.User)
	if err := u.UserRepository.FindByID(u.DB, user, userUUID, []string{"Roles", "Roles.Permissions"}); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("User not found: %s", userUUID.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrUserNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	return converter.UserToResponse(user), nil
}
