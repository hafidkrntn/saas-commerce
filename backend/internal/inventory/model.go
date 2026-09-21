package inventory

import (
	"time"

	"github.com/google/uuid"
)

// Warehouse represents the `warehouses` table.
type Warehouse struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  uuid.UUID `gorm:"column:tenant_id;type:uuid;not null;index"`
	Name      string    `gorm:"column:name;not null"`
	Address   string    `gorm:"column:address"`
	IsDefault bool      `gorm:"column:is_default;default:false"`
	CreatedBy string    `gorm:"column:created_by"`
	UpdatedBy string    `gorm:"column:updated_by"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

// StockMovement represents the `stock_movements` table.
type StockMovement struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    uuid.UUID  `gorm:"column:tenant_id;type:uuid;not null;index"`
	ProductID   uuid.UUID  `gorm:"column:product_id;type:uuid;not null;index"`
	WarehouseID *uuid.UUID `gorm:"column:warehouse_id;type:uuid"`
	Type        string     `gorm:"column:type;not null"` // in | out | adjustment | return
	Quantity    int        `gorm:"column:quantity;not null"`
	Reference   string     `gorm:"column:reference"`
	Note        string     `gorm:"column:note"`
	UserName    string     `gorm:"column:user_name"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (StockMovement) TableName() string {
	return "stock_movements"
}

// IncomingShipment represents the `incoming_shipments` table.
type IncomingShipment struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    uuid.UUID  `gorm:"column:tenant_id;type:uuid;not null;index"`
	ProductID   uuid.UUID  `gorm:"column:product_id;type:uuid;not null"`
	Quantity    int        `gorm:"column:quantity;not null"`
	ExpectedDate *time.Time `gorm:"column:expected_date"`
	Supplier    string     `gorm:"column:supplier"`
	Status      string     `gorm:"column:status;not null;default:'scheduled'"`
	Notes       string     `gorm:"column:notes"`
	CreatedBy   string     `gorm:"column:created_by"`
	UpdatedBy   string     `gorm:"column:updated_by"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (IncomingShipment) TableName() string {
	return "incoming_shipments"
}

// Shipment statuses
const (
	ShipmentScheduled = "scheduled"
	ShipmentInTransit = "in_transit"
	ShipmentDelivered = "delivered"
	ShipmentDelayed   = "delayed"
)

// =============================================================================
// Request DTOs
// =============================================================================

type CreateWarehouseRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=200"`
	Address   string `json:"address"`
	IsDefault bool   `json:"is_default"`
}

type UpdateWarehouseRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=2,max=200"`
	Address   *string `json:"address"`
	IsDefault *bool   `json:"is_default"`
}

type StockMovementRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=in out adjustment return"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	Reference string `json:"reference"`
	Note      string `json:"note"`
}

type CreateShipmentRequest struct {
	ProductID    uuid.UUID  `json:"product_id" binding:"required"`
	Quantity     int        `json:"quantity" binding:"required,gt=0"`
	ExpectedDate *time.Time `json:"expected_date"`
	Supplier     string     `json:"supplier"`
	Notes        string     `json:"notes"`
}

type UpdateShipmentRequest struct {
	Quantity     *int       `json:"quantity" binding:"omitempty,gt=0"`
	ExpectedDate *time.Time `json:"expected_date"`
	Supplier     *string    `json:"supplier"`
	Status       *string    `json:"status" binding:"omitempty,oneof=scheduled in_transit delivered delayed"`
	Notes        *string    `json:"notes"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type WarehouseResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

func (w *Warehouse) ToResponse() WarehouseResponse {
	return WarehouseResponse{
		ID:        w.ID,
		Name:      w.Name,
		Address:   w.Address,
		IsDefault: w.IsDefault,
		CreatedAt: w.CreatedAt,
	}
}

type StockMovementResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	Reference string    `json:"reference"`
	Note      string    `json:"note"`
	User      string    `json:"user"`
	Date      time.Time `json:"date"`
}

func (m *StockMovement) ToResponse() StockMovementResponse {
	return StockMovementResponse{
		ID:        m.ID,
		ProductID: m.ProductID,
		Type:      m.Type,
		Quantity:  m.Quantity,
		Reference: m.Reference,
		Note:      m.Note,
		User:      m.UserName,
		Date:      m.CreatedAt,
	}
}

type ShipmentResponse struct {
	ID           uuid.UUID  `json:"id"`
	ProductID    uuid.UUID  `json:"product_id"`
	Quantity     int        `json:"quantity"`
	ExpectedDate *time.Time `json:"expected_date"`
	Supplier     string     `json:"supplier"`
	Status       string     `json:"status"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (s *IncomingShipment) ToResponse() ShipmentResponse {
	return ShipmentResponse{
		ID:           s.ID,
		ProductID:    s.ProductID,
		Quantity:     s.Quantity,
		ExpectedDate: s.ExpectedDate,
		Supplier:     s.Supplier,
		Status:       s.Status,
		Notes:        s.Notes,
		CreatedAt:    s.CreatedAt,
	}
}
