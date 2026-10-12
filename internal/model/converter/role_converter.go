package converter

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
)

func RoleToResponse(role *entity.Role) *payload.RoleResponse {
	response := &payload.RoleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	for _, permission := range role.Permissions {
		response.Permissions = append(response.Permissions, payload.PermissionResponse{
			ID:          permission.ID,
			Code:        permission.Code,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
			CreatedAt:   permission.CreatedAt,
		})
	}

	return response
}
