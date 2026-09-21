package payment

import (
	"backend-go/internal/entities"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents the `payment_methods` table.
type PaymentMethod struct {
	ID          uuid.UUID       `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    uuid.UUID       `gorm:"column:tenant_id;type:uuid;not null;index"`
	Name        string          `gorm:"column:name;not null"`
	Type        string          `gorm:"column:type;not null"`
	IsEnabled   bool            `gorm:"column:is_enabled;default:true"`
	Description string          `gorm:"column:description"`
	Config      entities.JSONB    `gorm:"column:config;type:jsonb"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;autoUpdateTime"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}

// Payment represents the `payments` table.
type Payment struct {
	ID                   uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID             uuid.UUID  `gorm:"column:tenant_id;type:uuid;not null;index"`
	OrderID              uuid.UUID  `gorm:"column:order_id;type:uuid;not null;index"`
	PaymentMethodID      *uuid.UUID `gorm:"column:payment_method_id;type:uuid"`
	Provider             string     `gorm:"column:provider;not null;default:'mock'"`
	ProviderTransactionID string   `gorm:"column:provider_transaction_id"`
	Status               string     `gorm:"column:status;not null;default:'pending'"`
	Amount               float64    `gorm:"column:amount;not null"`
	Currency             string     `gorm:"column:currency;not null;default:'IDR'"`
	PaidAt               *time.Time `gorm:"column:paid_at"`
	CreatedAt            time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Payment) TableName() string {
	return "payments"
}

// Payment statuses
const (
	PayStatusPending  = "pending"
	PayStatusPaid     = "paid"
	PayStatusFailed   = "failed"
	PayStatusRefunded = "refunded"
)

// =============================================================================
// Request DTOs
// =============================================================================

type CreatePaymentMethodRequest struct {
	Name        string         `json:"name" binding:"required,min=2,max=100"`
	Type        string         `json:"type" binding:"required,oneof=card bank_transfer cod e_wallet"`
	IsEnabled   *bool          `json:"is_enabled"`
	Description string         `json:"description"`
	Config      map[string]any `json:"config"`
}

type UpdatePaymentMethodRequest struct {
	Name        *string        `json:"name" binding:"omitempty,min=2,max=100"`
	Type        *string        `json:"type" binding:"omitempty,oneof=card bank_transfer cod e_wallet"`
	IsEnabled   *bool          `json:"is_enabled"`
	Description *string        `json:"description"`
	Config      map[string]any `json:"config"`
}

type ChargeRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type PaymentMethodResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	IsEnabled   bool           `json:"is_enabled"`
	Description string         `json:"description"`
	Config      entities.JSONB    `json:"config,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (m *PaymentMethod) ToResponse() PaymentMethodResponse {
	return PaymentMethodResponse{
		ID:          m.ID,
		Name:        m.Name,
		Type:        m.Type,
		IsEnabled:   m.IsEnabled,
		Description: m.Description,
		Config:      m.Config,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type PaymentResponse struct {
	ID                     uuid.UUID  `json:"id"`
	OrderID                uuid.UUID  `json:"order_id"`
	PaymentMethodID        *uuid.UUID `json:"payment_method_id,omitempty"`
	Provider               string     `json:"provider"`
	ProviderTransactionID  string     `json:"provider_transaction_id"`
	Status                 string     `json:"status"`
	Amount                 float64    `json:"amount"`
	Currency               string     `json:"currency"`
	PaidAt                 *time.Time `json:"paid_at"`
	CreatedAt              time.Time  `json:"created_at"`
}

func (p *Payment) ToResponse() PaymentResponse {
	return PaymentResponse{
		ID:                    p.ID,
		OrderID:               p.OrderID,
		PaymentMethodID:       p.PaymentMethodID,
		Provider:              p.Provider,
		ProviderTransactionID: p.ProviderTransactionID,
		Status:                p.Status,
		Amount:                p.Amount,
		Currency:              p.Currency,
		PaidAt:                p.PaidAt,
		CreatedAt:             p.CreatedAt,
	}
}

type ChargeResponse struct {
	Payment     PaymentResponse `json:"payment"`
	RedirectURL string          `json:"redirect_url"`
	Status      string          `json:"status"`
}
