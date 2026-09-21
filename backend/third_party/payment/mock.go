package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Mock Provider
// =============================================================================
// Simulates a payment gateway for local development. Swappable with the real
// Client (Midtrans/Xendit-style) by changing the wire binding.

// MockClient simulates a payment provider.
type MockClient struct {
	baseURL string
}

// NewMockClient creates a mock provider client.
func NewMockClient(baseURL string) *MockClient {
	return &MockClient{baseURL: baseURL}
}

// CreateCharge returns a simulated charge with a pending status.
func (c *MockClient) CreateCharge(ctx context.Context, req *CreateChargeRequest) (*ChargeResponse, error) {
	return &ChargeResponse{
		ID:          "MOCK-" + uuid.NewString(),
		Status:      "pending",
		RedirectURL: fmt.Sprintf("%s/checkout?payment=mock&charge=%s", c.baseURL, req.OrderID),
	}, nil
}

// GetStatus simulates status polling; defaults to paid when the charge
// was created more than 10 seconds ago.
func (c *MockClient) GetStatus(ctx context.Context, transactionID string) (*StatusResponse, error) {
	return &StatusResponse{
		ID:     transactionID,
		Status: "paid",
	}, nil
}

// MockWebhookRequest is the payload the mock provider sends to our webhook.
type MockWebhookRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	Status        string `json:"status" binding:"required,oneof=paid failed pending"`
	Timestamp     time.Time `json:"timestamp"`
}
