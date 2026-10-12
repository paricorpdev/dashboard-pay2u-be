package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/constants"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/exception"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/logger"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/converter"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RBACUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	RoleRepository       *repository.RoleRepository
	PermissionRepository *repository.PermissionRepository
}

func NewRBACUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
) *RBACUseCase {
	return &RBACUseCase{
		DB:                   db,
		Log:                  log,
		RoleRepository:       roleRepo,
		PermissionRepository: permissionRepo,
	}
}

// FindAllRoles returns all available roles with their associated permissions.
func (u *RBACUseCase) FindAllRoles(ctx *fiber.Ctx) ([]payload.RoleResponse, error) {
	roles, err := u.RoleRepository.FindAll(u.DB)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to retrieve roles: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.RolesToResponses(roles), nil
}

// FindRoleByID returns detailed role info, user count, and permissions for a selected role.
func (u *RBACUseCase) FindRoleByID(ctx *fiber.Ctx, id uuid.UUID) (*payload.RoleDetailResponse, error) {
	var role entity.Role
	if err := u.RoleRepository.FindByID(u.DB, &role, id); err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Role not found: %s", id.String())
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrRoleNotFound, fiber.ErrNotFound),
		})
		ctx.Status(fiber.StatusNotFound)
		return nil, fiber.ErrNotFound
	}

	userCount, err := u.RoleRepository.CountUsersByRoleID(u.DB, id)
	if err != nil {
		logger.WithContext(ctx, u.Log).Warnf("Failed to count users for role %s: %+v", id.String(), err)
		userCount = 0
	}

	return converter.RoleToDetailResponse(&role, userCount), nil
}

// FindAllPermissions returns all registered system permissions.
func (u *RBACUseCase) FindAllPermissions(ctx *fiber.Ctx) ([]payload.PermissionResponse, error) {
	perms, err := u.PermissionRepository.FindAll(u.DB)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to retrieve permissions: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PermissionsToResponses(perms), nil
}

// FindPermissionsGrouped returns all registered system permissions grouped by resource.
func (u *RBACUseCase) FindPermissionsGrouped(ctx *fiber.Ctx) ([]payload.PermissionGroupResponse, error) {
	perms, err := u.PermissionRepository.FindAll(u.DB)
	if err != nil {
		logger.WithContext(ctx, u.Log).Errorf("Failed to retrieve permissions: %+v", err)
		ctx.Locals(constants.Response, payload.ResponseInterface{
			Status:       constants.ResponseNameError,
			ResponseData: exception.NewHandleError(ctx, constants.ErrInternalServerError, fiber.ErrInternalServerError),
		})
		ctx.Status(fiber.StatusInternalServerError)
		return nil, fiber.ErrInternalServerError
	}

	return converter.GroupPermissionsByResource(perms), nil
}
