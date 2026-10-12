package converter

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
)

func UserToResponse(user *entity.User) *payload.UserResponse {
	response := &payload.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	for _, role := range user.Roles {
		response.Roles = append(response.Roles, *RoleToResponse(role))
	}

	return response
}

func UserToPageResponse(page *payload.WebResponse[[]entity.User]) *payload.PageResponse[*payload.UserResponse] {
	response := new(payload.PageResponse[*payload.UserResponse])
	response.PageMetadata.Page = page.Paging.Page
	response.PageMetadata.PerPage = page.Paging.PerPage
	response.PageMetadata.TotalData = page.Paging.TotalData
	response.PageMetadata.TotalPage = page.Paging.TotalPage

	for _, user := range page.Data {
		response.Data = append(response.Data, UserToResponse(&user))
	}

	return response
}
