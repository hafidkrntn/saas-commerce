package payment

import (
	"backend-go/pkg/httpclient"
	"context"
	"time"
)

// =============================================================================
// Interface
// =============================================================================
// Define interface here so services can mock this dependency for testing.

type ClientInterface interface {
	CreateCharge(ctx context.Context, req *CreateChargeRequest) (*ChargeResponse, error)
	GetStatus(ctx context.Context, transactionID string) (*StatusResponse, error)
}

// =============================================================================
// Implementation
// =============================================================================

type Client struct {
	http *httpclient.Client
}

func NewClient(baseURL string, apiKey string) *Client {
	return &Client{
		http: httpclient.New(baseURL, httpclient.Options{
			Timeout: 15 * time.Second,
			Headers: map[string]string{
				"Authorization": "Bearer " + apiKey,
			},
		}),
	}
}

func (c *Client) CreateCharge(ctx context.Context, req *CreateChargeRequest) (*ChargeResponse, error) {
	var resp ChargeResponse
	if err := c.http.Post(ctx, "/v1/charges", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetStatus(ctx context.Context, transactionID string) (*StatusResponse, error) {
	var resp StatusResponse
	if err := c.http.Get(ctx, "/v1/transactions/"+transactionID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// =============================================================================
// Request / Response DTOs
// =============================================================================

type CreateChargeRequest struct {
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
}

type ChargeResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	RedirectURL string `json:"redirect_url"`
}

type StatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"` // pending, paid, failed, expired
}
