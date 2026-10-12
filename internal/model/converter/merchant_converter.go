package converter

import (
	"time"

	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
)

func MerchantToResponse(merchant *entity.Merchant) *payload.MerchantResponse {
	return &payload.MerchantResponse{
		ID:            merchant.ID.String(),
		Code:          merchant.Code,
		Name:          merchant.Name,
		ApiKey:        merchant.ApiKey,
		WebhookUrl:    merchant.WebhookUrl,
		WebhookSecret: merchant.WebhookSecret,
		IsActive:      merchant.IsActive,
		CreatedAt:     merchant.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     merchant.UpdatedAt.Format(time.RFC3339),
	}
}

func MerchantsToResponses(merchants []entity.Merchant) []payload.MerchantResponse {
	responses := make([]payload.MerchantResponse, 0, len(merchants))
	for i := range merchants {
		responses = append(responses, *MerchantToResponse(&merchants[i]))
	}
	return responses
}
