package payload

type QueryTransactionRequest struct {
	Page          int     `query:"page"`
	Limit         int     `query:"limit"`
	Name          string  `query:"name"`
	Phone         string  `query:"phone"`
	Email         string  `query:"email"`
	PaymentMethod string  `query:"payment_method"`
	Status        string  `query:"status"`
	MerchantReff  string  `query:"merchant_reff"`
	BillTitle     string  `query:"bill_title"`
	StatusPaid    *bool   `query:"status_paid"`
	MerchantID    *string `query:"merchant_id"`
	StartDate     string  `query:"start_date"`
	EndDate       string  `query:"end_date"`
	SortBy        string  `query:"sort_by"`
	SortDir       string  `query:"sort_dir"`
}

type TransactionResponse struct {
	ID              string           `json:"id"`
	RequestID       string           `json:"request_id"`
	Provider        string           `json:"provider"`
	ProviderToken   *string          `json:"provider_token"`
	MerchantReff    string           `json:"merchant_reff"`
	PaymentMethod   string           `json:"payment_method"`
	PaymentCode     *string          `json:"payment_code"`
	Status          string           `json:"status"`
	IsPaid          bool             `json:"is_paid"`
	Amount          int64            `json:"amount"`
	AmountAdmin     int64            `json:"amount_admin"`
	AmountDiscount  int64            `json:"amount_discount"`
	AmountTotal     int64            `json:"amount_total"`
	Currency        string           `json:"currency"`
	CustomerName    *string          `json:"customer_name"`
	CustomerPhone   *string          `json:"customer_phone"`
	CustomerEmail   *string          `json:"customer_email"`
	BillTitle       *string          `json:"bill_title"`
	BillDescription *string          `json:"bill_description"`
	PaymentReff     *string          `json:"payment_reff"`
	ExpiredAt       *string          `json:"expired_at"`
	PaidAt          *string          `json:"paid_at"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
	Merchant        *MerchantSummary `json:"merchant"`
}

type TransactionDetailResponse struct {
	TransactionResponse
	ProviderCallback  *string                  `json:"provider_callback"`
	WebhookDispatches []WebhookDispatchSummary `json:"webhook_dispatches"`
}

type MerchantSummary struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type WebhookDispatchSummary struct {
	ID          string  `json:"id"`
	TargetUrl   string  `json:"target_url"`
	EventType   string  `json:"event_type"`
	Status      string  `json:"status"`
	Attempts    int     `json:"attempts"`
	MaxAttempts int     `json:"max_attempts"`
	NextRetryAt *string `json:"next_retry_at"`
	CreatedAt   string  `json:"created_at"`
}
