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
		Permissions: make([]payload.PermissionResponse, 0),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	for _, permission := range role.Permissions {
		if permission != nil {
			response.Permissions = append(response.Permissions, *PermissionToResponse(permission))
		}
	}

	return response
}

func RolesToResponses(roles []entity.Role) []payload.RoleResponse {
	responses := make([]payload.RoleResponse, 0, len(roles))
	for i := range roles {
		responses = append(responses, *RoleToResponse(&roles[i]))
	}
	return responses
}

func RoleToListItemResponse(role *entity.Role) *payload.RoleListItemResponse {
	return &payload.RoleListItemResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func RolesToListResponses(roles []entity.Role) []payload.RoleListItemResponse {
	responses := make([]payload.RoleListItemResponse, 0, len(roles))
	for i := range roles {
		responses = append(responses, *RoleToListItemResponse(&roles[i]))
	}
	return responses
}

func RoleToDetailResponse(role *entity.Role, userCount int64) *payload.RoleDetailResponse {
	response := &payload.RoleDetailResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		UserCount:   userCount,
		Permissions: make([]payload.PermissionResponse, 0),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	for _, permission := range role.Permissions {
		if permission != nil {
			response.Permissions = append(response.Permissions, *PermissionToResponse(permission))
		}
	}

	return response
}

func PermissionToResponse(p *entity.Permission) *payload.PermissionResponse {
	return &payload.PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func PermissionsToResponses(perms []entity.Permission) []payload.PermissionResponse {
	responses := make([]payload.PermissionResponse, 0, len(perms))
	for i := range perms {
		responses = append(responses, *PermissionToResponse(&perms[i]))
	}
	return responses
}

func GroupPermissionsByResource(perms []entity.Permission) []payload.PermissionGroupResponse {
	groupMap := make(map[string][]payload.PermissionResponse)
	var order []string

	for i := range perms {
		r := perms[i].Resource
		if _, exists := groupMap[r]; !exists {
			order = append(order, r)
		}
		groupMap[r] = append(groupMap[r], *PermissionToResponse(&perms[i]))
	}

	result := make([]payload.PermissionGroupResponse, 0, len(order))
	for _, r := range order {
		result = append(result, payload.PermissionGroupResponse{
			Resource:    r,
			Permissions: groupMap[r],
		})
	}

	return result
}
