package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/hashing"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/keys"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/converter"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	RoleRepository *repository.RoleRepository
}

func NewUserUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepo,
		RoleRepository: roleRepo,
	}
}

func (u *UserUseCase) CreateAdmin(ctx *fiber.Ctx, request *payload.CreateAdminRequest) (*payload.UserResponse, error) {
	request.Username = strings.TrimSpace(request.Username)
	request.FullName = strings.TrimSpace(request.FullName)
	if request.Email != nil {
		email := strings.TrimSpace(*request.Email)
		if email == "" {
			request.Email = nil
		} else {
			request.Email = &email
		}
	}

	if err := u.Validate.Struct(request); err != nil {
		return nil, u.fail(ctx, fiber.StatusBadRequest, err.Error())
	}

	exists, err := u.UserRepository.ExistsByUsername(u.DB, request.Username)
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	if exists {
		return nil, u.fail(ctx, fiber.StatusConflict, constants.ErrUserAlreadyExists)
	}
	if request.Email != nil {
		exists, err = u.UserRepository.ExistsByEmail(u.DB, *request.Email)
		if err != nil {
			return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
		}
		if exists {
			return nil, u.fail(ctx, fiber.StatusConflict, constants.ErrEmailAlreadyExists)
		}
	}

	passwordHash, err := hashing.HashPassword(request.Password)
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	msMasterKey, err := keys.GenerateDashboardMSKey()
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}

	roles, err := u.RoleRepository.FindByCodes(u.DB, []string{constants.RoleAdmin})
	if err != nil || len(roles) != 1 {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrRoleNotFound)
	}

	now := time.Now().UTC()
	user := &entity.User{
		ID:           uuid.New(),
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: passwordHash,
		MSMasterKey:  &msMasterKey,
		FullName:     request.FullName,
		Status:       constants.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tx := u.DB.WithContext(ctx.Context()).Begin()
	if tx.Error != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	defer tx.Rollback()

	if err := u.UserRepository.Create(tx, user); err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	if err := tx.Model(user).Association("Roles").Append(&roles[0]); err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	if err := tx.Preload("Roles.Permissions").First(user, "id = ?", user.ID).Error; err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) FindAdmins(ctx *fiber.Ctx, request *payload.QueryAdminRequest) (*payload.PageResponse[payload.UserResponse], error) {
	if request.Status != "" {
		request.Status = strings.ToUpper(strings.TrimSpace(request.Status))
		if request.Status != constants.StatusActive && request.Status != constants.StatusInactive {
			return nil, u.fail(ctx, fiber.StatusBadRequest, "status must be ACTIVE or INACTIVE")
		}
	}

	users, metadata, err := u.UserRepository.FindAdmins(u.DB, request)
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	response := &payload.PageResponse[payload.UserResponse]{
		Data:         make([]payload.UserResponse, 0, len(users)),
		PageMetadata: *metadata,
	}
	for i := range users {
		response.Data = append(response.Data, *converter.UserToResponse(&users[i]))
	}
	return response, nil
}

func (u *UserUseCase) FindAdmin(ctx *fiber.Ctx, id uuid.UUID) (*payload.UserResponse, error) {
	user, err := u.UserRepository.FindAdminByID(u.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, u.fail(ctx, fiber.StatusNotFound, constants.ErrUserNotFound)
	}
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) UpdateAdmin(ctx *fiber.Ctx, id uuid.UUID, request *payload.UpdateAdminRequest) (*payload.UserResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, u.fail(ctx, fiber.StatusBadRequest, err.Error())
	}
	user, err := u.UserRepository.FindAdminByID(u.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, u.fail(ctx, fiber.StatusNotFound, constants.ErrUserNotFound)
	}
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}

	updates := make(map[string]any)
	if request.Username != nil {
		value := strings.TrimSpace(*request.Username)
		if value == "" {
			return nil, u.fail(ctx, fiber.StatusBadRequest, "username cannot be empty")
		}
		exists, err := u.UserRepository.ValidateExistByUsername(u.DB, value, id)
		if err != nil {
			return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
		}
		if exists {
			return nil, u.fail(ctx, fiber.StatusConflict, constants.ErrUserAlreadyExists)
		}
		updates["username"] = value
	}
	if request.Email != nil {
		value := strings.TrimSpace(*request.Email)
		if value == "" {
			updates["email"] = nil
		} else {
			exists, err := u.UserRepository.ValidateExistByEmail(u.DB, value, id)
			if err != nil {
				return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
			}
			if exists {
				return nil, u.fail(ctx, fiber.StatusConflict, constants.ErrEmailAlreadyExists)
			}
			updates["email"] = value
		}
	}
	if request.FullName != nil {
		value := strings.TrimSpace(*request.FullName)
		if value == "" {
			return nil, u.fail(ctx, fiber.StatusBadRequest, "full_name cannot be empty")
		}
		updates["full_name"] = value
	}
	if request.Password != nil && *request.Password != "" {
		passwordHash, err := hashing.HashPassword(*request.Password)
		if err != nil {
			return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
		}
		updates["password_hash"] = passwordHash
	}
	if request.Status != nil {
		status := strings.ToUpper(strings.TrimSpace(*request.Status))
		if status != constants.StatusActive && status != constants.StatusInactive {
			return nil, u.fail(ctx, fiber.StatusBadRequest, "status must be ACTIVE or INACTIVE")
		}
		updates["status"] = status
	}
	if len(updates) == 0 {
		return nil, u.fail(ctx, fiber.StatusBadRequest, "at least one field must be provided")
	}
	updates["updated_at"] = time.Now().UTC()
	if err := u.DB.Model(user).Updates(updates).Error; err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	if err := u.DB.Preload("Roles.Permissions").First(user, "id = ?", id).Error; err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) DeleteAdmin(ctx *fiber.Ctx, id uuid.UUID) (*payload.UserResponse, error) {
	user, err := u.UserRepository.FindAdminByID(u.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, u.fail(ctx, fiber.StatusNotFound, constants.ErrUserNotFound)
	}
	if err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	updatedAt := time.Now().UTC()
	if err := u.DB.Model(user).Updates(map[string]any{
		"status":     constants.StatusInactive,
		"updated_at": updatedAt,
	}).Error; err != nil {
		return nil, u.fail(ctx, fiber.StatusInternalServerError, constants.ErrInternalServerError)
	}
	user.Status = constants.StatusInactive
	user.UpdatedAt = updatedAt
	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) fail(ctx *fiber.Ctx, status int, message string) error {
	fiberErr := fiber.NewError(status, message)
	ctx.Locals(constants.Response, payload.ResponseInterface{
		Status:       constants.ResponseNameError,
		ResponseData: exception.NewHandleError(ctx, message, fiberErr),
	})
	ctx.Status(status)
	return fiberErr
}
