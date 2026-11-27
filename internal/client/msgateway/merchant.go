package msgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
)

type MerchantClient struct {
	baseURL string
	http    *http.Client
}

type MerchantListQuery struct {
	Page     int
	PerPage  int
	Search   string
	IsActive *bool
	SortBy   string
	SortDir  string
}

type DisbursementListQuery struct {
	Page       int
	PerPage    int
	Search     string
	Status     string
	MerchantID string
	Month      int
	Year       int
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("ms merchant API returned status %d: %s", e.StatusCode, e.Message)
}

type envelope[T any] struct {
	ResponseData T `json:"response_data"`
	Error        struct {
		Message string `json:"message"`
	} `json:"error"`
}

type MerchantPage struct {
	Data   []payload.MerchantResponse `json:"data"`
	Paging payload.PageMetadata       `json:"paging"`
}

func NewMerchantClient(baseURL string) *MerchantClient {
	return &MerchantClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *MerchantClient) Create(ctx context.Context, masterKey string, request payload.CreateMerchantRequest) (*payload.MerchantResponse, error) {
	body := map[string]any{
		"code":        request.Code,
		"name":        request.Name,
		"webhook_url": stringValue(request.WebhookUrl),
	}
	var merchant payload.MerchantResponse
	if err := c.do(ctx, masterKey, http.MethodPost, "/api/v1/admin/merchants/", nil, body, &merchant); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (c *MerchantClient) List(ctx context.Context, masterKey string, query MerchantListQuery) (*MerchantPage, error) {
	values := url.Values{}
	values.Set("page", strconv.Itoa(query.Page))
	values.Set("per_page", strconv.Itoa(query.PerPage))
	if query.Search != "" {
		values.Set("search", query.Search)
	}
	if query.IsActive != nil {
		values.Set("is_active", strconv.FormatBool(*query.IsActive))
	}
	if query.SortBy != "" {
		values.Set("sort_by", query.SortBy)
	}
	if query.SortDir != "" {
		values.Set("sort_dir", query.SortDir)
	}

	var response MerchantPage
	if err := c.do(ctx, masterKey, http.MethodGet, "/api/v1/admin/merchants/", values, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *MerchantClient) Get(ctx context.Context, masterKey, id string) (*payload.MerchantResponse, error) {
	var merchant payload.MerchantResponse
	if err := c.do(ctx, masterKey, http.MethodGet, "/api/v1/admin/merchants/"+url.PathEscape(id), nil, nil, &merchant); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (c *MerchantClient) Update(ctx context.Context, masterKey, id string, request map[string]any) (*payload.MerchantResponse, error) {
	var merchant payload.MerchantResponse
	path := "/api/v1/admin/merchants/" + url.PathEscape(id)
	if err := c.do(ctx, masterKey, http.MethodPut, path, nil, request, &merchant); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (c *MerchantClient) RotateAPIKey(ctx context.Context, masterKey, id string) (*payload.MerchantResponse, error) {
	var merchant payload.MerchantResponse
	path := "/api/v1/admin/merchants/" + url.PathEscape(id) + "/rotate-key"
	if err := c.do(ctx, masterKey, http.MethodPost, path, nil, nil, &merchant); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (c *MerchantClient) RotateWebhookSecret(ctx context.Context, masterKey, id string) (*payload.MerchantResponse, error) {
	var merchant payload.MerchantResponse
	path := "/api/v1/admin/merchants/" + url.PathEscape(id) + "/rotate-webhook-secret"
	if err := c.do(ctx, masterKey, http.MethodPost, path, nil, nil, &merchant); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (c *MerchantClient) ListDisbursements(ctx context.Context, masterKey string, query DisbursementListQuery) (*payload.PageResponse[payload.DisbursementResponse], error) {
	values := url.Values{}
	values.Set("page", strconv.Itoa(query.Page))
	values.Set("per_page", strconv.Itoa(query.PerPage))
	if query.Search != "" {
		values.Set("search", query.Search)
	}
	if query.Status != "" {
		values.Set("status", query.Status)
	}
	if query.MerchantID != "" {
		values.Set("merchant_id", query.MerchantID)
	}
	if query.Month > 0 {
		values.Set("month", strconv.Itoa(query.Month))
	}
	if query.Year > 0 {
		values.Set("year", strconv.Itoa(query.Year))
	}
	var response payload.PageResponse[payload.DisbursementResponse]
	if err := c.do(ctx, masterKey, http.MethodGet, "/api/v1/admin/disbursements/", values, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *MerchantClient) GetDisbursement(ctx context.Context, masterKey, id string) (*payload.DisbursementDetailResponse, error) {
	var response payload.DisbursementDetailResponse
	path := "/api/v1/admin/disbursements/" + url.PathEscape(id)
	if err := c.do(ctx, masterKey, http.MethodGet, path, nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *MerchantClient) do(ctx context.Context, masterKey, method, path string, query url.Values, body any, output any) error {
	if c.baseURL == "" || strings.TrimSpace(masterKey) == "" {
		return errors.New("ms_gateway.base_url and the signed-in user's MS key must be configured")
	}

	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode ms merchant request: %w", err)
		}
		bodyReader = bytes.NewReader(encoded)
	}

	request, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("build ms merchant request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-API-Key", masterKey)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("call ms merchant service: %w", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(io.LimitReader(response.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("read ms merchant response: %w", err)
	}

	var result envelope[json.RawMessage]
	if len(data) > 0 {
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("decode ms merchant response: %w", err)
		}
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		message := result.Error.Message
		if message == "" {
			message = http.StatusText(response.StatusCode)
		}
		return &APIError{StatusCode: response.StatusCode, Message: message}
	}
	if len(result.ResponseData) == 0 || string(result.ResponseData) == "null" {
		return errors.New("ms merchant response_data is empty")
	}
	if err := json.Unmarshal(result.ResponseData, output); err != nil {
		return fmt.Errorf("decode ms merchant response_data: %w", err)
	}
	return nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
