package payload

type PageRequest struct {
	Page    int                    `json:"page" validate:"gte=1"`
	PerPage int                    `json:"per_page" validate:"gte=1,lte=100"`
	Sort    map[string]interface{} `json:"sort"`
	Include []string               `json:"include"`
}

type WebResponse[T any] struct {
	Data   T             `json:"data"`
	Paging *PageMetadata `json:"paging,omitempty"`
	Errors string        `json:"errors,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data"`
	PageMetadata PageMetadata `json:"paging"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
}

type ResponseInterface struct {
	Status       string `json:"status"`
	ResponseData any    `json:"response_data"`
}

type Response struct {
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	ResponseData any    `json:"response_data"`
}

type ResponseError struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Error      any    `json:"error"`
	RequestId  string `json:"request_id"`
}
