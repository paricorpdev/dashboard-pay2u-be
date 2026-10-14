package payload

type MerchantOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MerchantResponse struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	ApiKey        string  `json:"api_key"`
	WebhookUrl    *string `json:"webhook_url"`
	WebhookSecret *string `json:"webhook_secret"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CreateMerchantRequest struct {
	Code       string  `json:"code" validate:"required,min=3,max=50"`
	Name       string  `json:"name" validate:"required,min=3,max=100"`
	WebhookUrl *string `json:"webhook_url" validate:"omitempty,url"`
}

type UpdateMerchantRequest struct {
	Name       string  `json:"name" validate:"required,min=3,max=100"`
	WebhookUrl *string `json:"webhook_url" validate:"omitempty,url"`
	IsActive   *bool   `json:"is_active" validate:"required"`
}

type RegenerateMerchantSecretRequest struct {
	RotateApiKey        bool `json:"rotate_api_key"`
	RotateWebhookSecret bool `json:"rotate_webhook_secret"`
}

type QueryMerchantRequest struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Search   string `query:"search"`
	IsActive *bool  `query:"is_active"`
	SortBy   string `query:"sort_by"`
	SortDir  string `query:"sort_dir"`
}
