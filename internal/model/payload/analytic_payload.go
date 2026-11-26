package payload

type AnalyticFilterRequest struct {
	Month      int     `query:"month"`
	Year       int     `query:"year"`
	MerchantID *string `query:"merchant_id"`
}

type AnalyticSummaryResponse struct {
	Month                 int     `json:"month"`
	Year                  int     `json:"year"`
	TotalVolumeAmount     int64   `json:"total_volume_amount"`
	TotalTransactionCount int64   `json:"total_transaction_count"`
	SuccessCount          int64   `json:"success_count"`
	PendingCount          int64   `json:"pending_count"`
	ExpiredCount          int64   `json:"expired_count"`
	SuccessRate           float64 `json:"success_rate"`
	TotalAdminFee         int64   `json:"total_admin_fee"`
}

type DailyTrendItem struct {
	Date         string `json:"date"`
	Day          int    `json:"day"`
	TotalAmount  int64  `json:"total_amount"`
	TotalCount   int64  `json:"total_count"`
	SuccessCount int64  `json:"success_count"`
	PendingCount int64  `json:"pending_count"`
	ExpiredCount int64  `json:"expired_count"`
}

type PaymentMethodDistributionItem struct {
	PaymentMethod string  `json:"payment_method"`
	Count         int64   `json:"count"`
	Amount        int64   `json:"amount"`
	Percentage    float64 `json:"percentage"`
}

type StatusDistributionItem struct {
	Status     string  `json:"status"`
	Count      int64   `json:"count"`
	Amount     int64   `json:"amount"`
	Percentage float64 `json:"percentage"`
}

type DisbursementSummaryResponse struct {
	Month                   int     `json:"month"`
	Year                    int     `json:"year"`
	TotalVolumeAmount       int64   `json:"total_volume_amount"`
	TotalDisbursementCount  int64   `json:"total_disbursement_count"`
	SuccessCount            int64   `json:"success_count"`
	WaitingSettlementCount  int64   `json:"waiting_settlement_count"`
	WaitingSettlementAmount int64   `json:"waiting_settlement_amount"`
	ProcessingCount         int64   `json:"processing_count"`
	ProcessingAmount        int64   `json:"processing_amount"`
	FailedCount             int64   `json:"failed_count"`
	FailedAmount            int64   `json:"failed_amount"`
	CancelledCount          int64   `json:"cancelled_count"`
	CancelledAmount         int64   `json:"cancelled_amount"`
	SuccessRate             float64 `json:"success_rate"`
	TotalAdminFee           int64   `json:"total_admin_fee"`
}

type DisbursementDailyTrendItem struct {
	Date                   string `json:"date"`
	Day                    int    `json:"day"`
	TotalAmount            int64  `json:"total_amount"`
	TotalCount             int64  `json:"total_count"`
	SuccessCount           int64  `json:"success_count"`
	WaitingSettlementCount int64  `json:"waiting_settlement_count"`
	ProcessingCount        int64  `json:"processing_count"`
	FailedCount            int64  `json:"failed_count"`
	CancelledCount         int64  `json:"cancelled_count"`
}

type DisbursementBankDistributionItem struct {
	BankName   string  `json:"bank_name"`
	Count      int64   `json:"count"`
	Amount     int64   `json:"amount"`
	Percentage float64 `json:"percentage"`
}
