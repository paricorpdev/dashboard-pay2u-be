package payload

import "time"

type QueryDisbursementRequest struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Search     string `query:"search"`
	Status     string `query:"status"`
	MerchantID string `query:"merchant_id"`
	Month      int    `query:"month"`
	Year       int    `query:"year"`
}

type DisbursementMerchantSummary struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type DisbursementResponse struct {
	ID                      string                       `json:"id"`
	TransactionID           string                       `json:"transaction_id"`
	RequestID               string                       `json:"request_id"`
	MerchantRef             string                       `json:"merchant_ref"`
	Provider                string                       `json:"provider"`
	ProviderDisbursementRef *string                      `json:"provider_disbursement_ref"`
	ChannelCode             string                       `json:"channel_code"`
	Amount                  int64                        `json:"amount"`
	AmountAdmin             int64                        `json:"amount_admin"`
	AmountDiscount          int64                        `json:"amount_discount"`
	AmountTotal             int64                        `json:"amount_total"`
	Currency                string                       `json:"currency"`
	BeneficiaryName         string                       `json:"beneficiary_name"`
	BeneficiaryBankCode     string                       `json:"beneficiary_bank_code"`
	BeneficiaryBankName     string                       `json:"beneficiary_bank_name"`
	Note                    string                       `json:"note"`
	Status                  string                       `json:"status"`
	EligibleAt              *time.Time                   `json:"eligible_at"`
	InquiryAt               *time.Time                   `json:"inquiry_at"`
	TransferredAt           *time.Time                   `json:"transferred_at"`
	Attempts                int                          `json:"attempts"`
	MaxAttempts             int                          `json:"max_attempts"`
	PaidAt                  *time.Time                   `json:"paid_at"`
	CreatedAt               time.Time                    `json:"created_at"`
	UpdatedAt               time.Time                    `json:"updated_at"`
	Merchant                *DisbursementMerchantSummary `json:"merchant"`
}

type DisbursementDetailResponse struct {
	DisbursementResponse
	BeneficiaryAccountNumber string `json:"beneficiary_account_number"`
}
