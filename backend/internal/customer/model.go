package customer

import (
	"backend-go/internal/entities"
	"time"

	"github.com/google/uuid"
)

// Customer represents the `customers` table.
type Customer struct {
	ID                uuid.UUID           `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID          uuid.UUID           `gorm:"column:tenant_id;type:uuid;not null;index"`
	Name              string              `gorm:"column:name;not null"`
	Email             string              `gorm:"column:email"`
	Phone             string              `gorm:"column:phone"`
	Avatar            string              `gorm:"column:avatar"`
	Status            string              `gorm:"column:status;not null;default:'active'"`
	Tags              entities.StringList `gorm:"column:tags;type:jsonb"`
	TotalOrders       int                 `gorm:"column:total_orders;not null;default:0"`
	TotalSpent        float64             `gorm:"column:total_spent;not null;default:0"`
	LifetimeValue     float64             `gorm:"column:lifetime_value;not null;default:0"`
	AverageOrderValue float64             `gorm:"column:average_order_value;not null;default:0"`
	LastOrderAt       *time.Time          `gorm:"column:last_order_at"`
	Addresses         []CustomerAddress   `gorm:"foreignKey:CustomerID"`
	CreatedBy         string              `gorm:"column:created_by"`
	UpdatedBy         string              `gorm:"column:updated_by"`
	CreatedAt         time.Time           `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time           `gorm:"column:updated_at;autoUpdateTime"`
}

func (Customer) TableName() string {
	return "customers"
}

// CustomerAddress represents the `customer_addresses` table.
type CustomerAddress struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID   uuid.UUID `gorm:"column:tenant_id;type:uuid;not null"`
	CustomerID uuid.UUID `gorm:"column:customer_id;type:uuid;not null;index"`
	Label      string    `gorm:"column:label;default:'home'"`
	Line1      string    `gorm:"column:line1;not null"`
	Line2      string    `gorm:"column:line2"`
	City       string    `gorm:"column:city;not null"`
	State      string    `gorm:"column:state;not null"`
	Zip        string    `gorm:"column:zip"`
	Country    string    `gorm:"column:country;not null;default:'Indonesia'"`
	IsDefault  bool      `gorm:"column:is_default;default:false"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CustomerAddress) TableName() string {
	return "customer_addresses"
}

// CustomerActivity represents the `customer_activities` table.
type CustomerActivity struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID   uuid.UUID `gorm:"column:tenant_id;type:uuid;not null"`
	CustomerID uuid.UUID `gorm:"column:customer_id;type:uuid;not null;index"`
	Type       string    `gorm:"column:type;not null"` // order | login | review | support | payment
	Description string   `gorm:"column:description;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (CustomerActivity) TableName() string {
	return "customer_activities"
}

// Customer statuses
const (
	StatusActive  = "active"
	StatusInactive = "inactive"
	StatusBlocked = "blocked"
)

// =============================================================================
// Request DTOs
// =============================================================================

type CreateCustomerRequest struct {
	Name     string   `json:"name" binding:"required,min=2,max=200"`
	Email    string   `json:"email" binding:"omitempty,email"`
	Phone    string   `json:"phone"`
	Avatar   string   `json:"avatar"`
	Status   string   `json:"status" binding:"omitempty,oneof=active inactive blocked"`
	Tags     []string `json:"tags"`
	Address  *AddressRequest `json:"address"`
}

type UpdateCustomerRequest struct {
	Name    *string   `json:"name" binding:"omitempty,min=2,max=200"`
	Email   *string   `json:"email" binding:"omitempty,email"`
	Phone   *string   `json:"phone"`
	Avatar  *string   `json:"avatar"`
	Status  *string   `json:"status" binding:"omitempty,oneof=active inactive blocked"`
	Tags    []string  `json:"tags"`
}

type AddressRequest struct {
	Label     string `json:"label"`
	Line1     string `json:"line1" binding:"required"`
	Line2     string `json:"line2"`
	City      string `json:"city" binding:"required"`
	State     string `json:"state" binding:"required"`
	Zip       string `json:"zip"`
	Country   string `json:"country"`
	IsDefault bool   `json:"is_default"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type AddressResponse struct {
	ID        uuid.UUID `json:"id"`
	Label     string    `json:"label"`
	Line1     string    `json:"line1"`
	Line2     string    `json:"line2"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Zip       string    `json:"zip"`
	Country   string    `json:"country"`
	IsDefault bool      `json:"is_default"`
}

type ActivityResponse struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type CustomerResponse struct {
	ID                uuid.UUID          `json:"id"`
	Name              string             `json:"name"`
	Email             string             `json:"email"`
	Phone             string             `json:"phone"`
	Avatar            string             `json:"avatar"`
	Status            string             `json:"status"`
	Tags              []string           `json:"tags"`
	TotalOrders       int                `json:"total_orders"`
	TotalSpent        float64            `json:"total_spent"`
	LifetimeValue     float64            `json:"lifetime_value"`
	AverageOrderValue float64            `json:"average_order_value"`
	LastOrderDate     *time.Time         `json:"last_order_date"`
	Addresses         []AddressResponse  `json:"addresses"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

func (c *Customer) ToResponse() CustomerResponse {
	addresses := make([]AddressResponse, len(c.Addresses))
	for i, a := range c.Addresses {
		addresses[i] = AddressResponse{
			ID:        a.ID,
			Label:     a.Label,
			Line1:     a.Line1,
			Line2:     a.Line2,
			City:      a.City,
			State:     a.State,
			Zip:       a.Zip,
			Country:   a.Country,
			IsDefault: a.IsDefault,
		}
	}

	return CustomerResponse{
		ID:                c.ID,
		Name:              c.Name,
		Email:             c.Email,
		Phone:             c.Phone,
		Avatar:            c.Avatar,
		Status:            c.Status,
		Tags:              c.Tags,
		TotalOrders:       c.TotalOrders,
		TotalSpent:        c.TotalSpent,
		LifetimeValue:     c.LifetimeValue,
		AverageOrderValue: c.AverageOrderValue,
		LastOrderDate:     c.LastOrderAt,
		Addresses:         addresses,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}
