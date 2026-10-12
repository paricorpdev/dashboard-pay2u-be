package converter

import (
	"time"

	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
)

func TransactionToResponse(t *entity.Transaction) *payload.TransactionResponse {
	var expiredAtStr *string
	if t.ExpiredAt != nil {
		s := t.ExpiredAt.Format(time.RFC3339)
		expiredAtStr = &s
	}

	var paidAtStr *string
	if t.PaidAt != nil {
		s := t.PaidAt.Format(time.RFC3339)
		paidAtStr = &s
	}

	var merchantSummary *payload.MerchantSummary
	if t.Merchant != nil {
		merchantSummary = &payload.MerchantSummary{
			ID:   t.Merchant.ID.String(),
			Code: t.Merchant.Code,
			Name: t.Merchant.Name,
		}
	}

	return &payload.TransactionResponse{
		ID:              t.ID.String(),
		RequestID:       t.RequestID,
		Provider:        t.Provider,
		ProviderToken:   t.ProviderToken,
		MerchantReff:    t.MerchantReff,
		PaymentMethod:   t.PaymentMethod,
		PaymentCode:     t.PaymentCode,
		Status:          t.Status,
		IsPaid:          t.PaidAt != nil,
		Amount:          t.Amount,
		AmountAdmin:     t.AmountAdmin,
		AmountDiscount:  t.AmountDiscount,
		AmountTotal:     t.AmountTotal,
		Currency:        t.Currency,
		CustomerName:    t.CustomerName,
		CustomerPhone:   t.CustomerPhone,
		CustomerEmail:   t.CustomerEmail,
		BillTitle:       t.BillTitle,
		BillDescription: t.BillDescription,
		PaymentReff:     t.PaymentReff,
		ExpiredAt:       expiredAtStr,
		PaidAt:          paidAtStr,
		CreatedAt:       t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       t.UpdatedAt.Format(time.RFC3339),
		Merchant:        merchantSummary,
	}
}

func TransactionsToResponses(transactions []entity.Transaction) []payload.TransactionResponse {
	responses := make([]payload.TransactionResponse, 0, len(transactions))
	for i := range transactions {
		item := TransactionToResponse(&transactions[i])
		// Optimization: Omit long payment codes (such as 250+ char EMVCo QRIS strings) in list view to reduce bandwidth bloat.
		// Detailed full payment code is retained in GET /api/v1/transactions/:id.
		if item.PaymentCode != nil && len(*item.PaymentCode) > 32 {
			item.PaymentCode = nil
		}
		responses = append(responses, *item)
	}
	return responses
}

func TransactionToDetailResponse(t *entity.Transaction) *payload.TransactionDetailResponse {
	base := TransactionToResponse(t)

	dispatches := make([]payload.WebhookDispatchSummary, 0, len(t.WebhookDispatches))
	for _, d := range t.WebhookDispatches {
		var nextRetry *string
		if d.NextRetryAt != nil {
			s := d.NextRetryAt.Format(time.RFC3339)
			nextRetry = &s
		}
		dispatches = append(dispatches, payload.WebhookDispatchSummary{
			ID:          d.ID.String(),
			TargetUrl:   d.TargetUrl,
			EventType:   d.EventType,
			Status:      d.Status,
			Attempts:    d.Attempts,
			MaxAttempts: d.MaxAttempts,
			NextRetryAt: nextRetry,
			CreatedAt:   d.CreatedAt.Format(time.RFC3339),
		})
	}

	return &payload.TransactionDetailResponse{
		TransactionResponse: *base,
		ProviderCallback:    t.ProviderCallback,
		WebhookDispatches:   dispatches,
	}
}
