package order

import (
	"backend-go/internal/entities"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Database Models
// =============================================================================

// Order represents the `orders` table.
type Order struct {
	ID               uuid.UUID          `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID         uuid.UUID          `gorm:"column:tenant_id;type:uuid;not null;index"`
	OrderNumber      string             `gorm:"column:order_number;not null"`
	CustomerID       *uuid.UUID         `gorm:"column:customer_id;type:uuid"`
	CustomerName     string             `gorm:"column:customer_name;not null"`
	Status           OrderStatus        `gorm:"column:status;not null;default:'pending'"`
	PaymentStatus    PaymentStatus      `gorm:"column:payment_status;not null;default:'pending'"`
	PaymentMethodID  *uuid.UUID         `gorm:"column:payment_method_id;type:uuid"`
	ShippingMethod   string             `gorm:"column:shipping_method"`
	ShippingAddress  entities.JSONB     `gorm:"column:shipping_address;type:jsonb"`
	Subtotal         float64            `gorm:"column:subtotal;not null;default:0"`
	ShippingCost     float64            `gorm:"column:shipping_cost;not null;default:0"`
	Tax              float64            `gorm:"column:tax;not null;default:0"`
	Discount         float64            `gorm:"column:discount;not null;default:0"`
	TotalAmount      float64            `gorm:"column:total_amount;not null;default:0"`
	Notes            string             `gorm:"column:notes"`
	Items            []OrderItem        `gorm:"foreignKey:OrderID"`
	Timeline         []OrderTimeline    `gorm:"foreignKey:OrderID"`
	CreatedBy        string             `gorm:"column:created_by"`
	UpdatedBy        string             `gorm:"column:updated_by"`
	CreatedAt        time.Time          `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time          `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        *time.Time         `gorm:"column:deleted_at"`

	// Transient fields (joined, not stored)
	CustomerEmail     string `gorm:"-"`
	PaymentMethodName string `gorm:"-"`
}

func (Order) TableName() string {
	return "orders"
}

// OrderItem represents the `order_items` table.
type OrderItem struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     uuid.UUID `gorm:"column:tenant_id;type:uuid;not null"`
	OrderID      uuid.UUID `gorm:"column:order_id;type:uuid;not null;index"`
	ProductID    uuid.UUID `gorm:"column:product_id;type:uuid;not null"`
	ProductName  string    `gorm:"column:product_name;not null"`
	ProductSKU   string    `gorm:"column:product_sku;not null"`
	ProductImage string    `gorm:"column:product_image"`
	Quantity     int       `gorm:"column:quantity;not null"`
	UnitPrice    float64   `gorm:"column:unit_price;not null"`
	Subtotal     float64   `gorm:"column:subtotal;not null"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

// OrderTimeline represents the `order_timeline` table.
type OrderTimeline struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    uuid.UUID `gorm:"column:tenant_id;type:uuid;not null"`
	OrderID     uuid.UUID `gorm:"column:order_id;type:uuid;not null;index"`
	Type        string    `gorm:"column:type;not null"`
	Title       string    `gorm:"column:title;not null"`
	Description string    `gorm:"column:description"`
	UserName    string    `gorm:"column:user_name"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (OrderTimeline) TableName() string {
	return "order_timeline"
}

// Order statuses (frontend-compatible)
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
	StatusRefunded  OrderStatus = "refunded"
)

// Payment statuses
type PaymentStatus string

const (
	PaymentPending           PaymentStatus = "pending"
	PaymentPaid              PaymentStatus = "paid"
	PaymentFailed            PaymentStatus = "failed"
	PaymentRefunded          PaymentStatus = "refunded"
	PaymentPartiallyRefunded PaymentStatus = "partially_refunded"
)

// validTransitions maps an order status to the statuses it may move to.
var validTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:   {StatusConfirmed, StatusCancelled},
	StatusConfirmed: {StatusProcessing, StatusCancelled},
	StatusProcessing: {StatusShipped, StatusCancelled},
	StatusShipped:   {StatusDelivered, StatusCancelled},
	StatusDelivered: {StatusRefunded},
	StatusCancelled: {},
	StatusRefunded:  {},
}

// CanTransition reports whether from -> to is a valid status move.
func CanTransition(from, to OrderStatus) bool {
	for _, next := range validTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

// =============================================================================
// Request DTOs
// =============================================================================

type CreateOrderRequest struct {
	CustomerID      *uuid.UUID          `json:"customer_id"`
	CustomerName    string              `json:"customer_name" binding:"required"`
	CustomerEmail   string              `json:"customer_email"`
	Items           []OrderItemRequest  `json:"items" binding:"required,min=1,dive"`
	PaymentMethodID *uuid.UUID          `json:"payment_method_id"`
	ShippingMethod  string              `json:"shipping_method"`
	ShippingAddress entities.JSONB      `json:"shipping_address"`
	ShippingCost    float64             `json:"shipping_cost"`
	Discount        float64             `json:"discount"`
	Notes           string              `json:"notes"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=confirmed processing shipped delivered cancelled refunded"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type OrderResponse struct {
	ID               uuid.UUID              `json:"id"`
	OrderNumber      string                 `json:"order_number"`
	Customer         CustomerBrief          `json:"customer"`
	Status           OrderStatus            `json:"status"`
	PaymentStatus    PaymentStatus          `json:"payment_status"`
	PaymentMethodID  *uuid.UUID             `json:"payment_method_id,omitempty"`
	PaymentMethod    string                 `json:"payment_method"`
	ShippingMethod   string                 `json:"shipping_method"`
	ShippingAddress  entities.JSONB         `json:"shipping_address,omitempty"`
	Subtotal         float64                `json:"subtotal"`
	ShippingCost     float64                `json:"shipping_cost"`
	Tax              float64                `json:"tax"`
	Discount         float64                `json:"discount"`
	TotalAmount      float64                `json:"total"`
	Notes            string                 `json:"notes"`
	Items            []OrderItemResponse    `json:"items"`
	Timeline         []TimelineResponse     `json:"timeline"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// CustomerBrief is the lightweight customer reference attached to an order.
type CustomerBrief struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type OrderItemResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductSKU   string    `json:"product_sku"`
	ProductImage string    `json:"product_image"`
	Quantity     int       `json:"quantity"`
	UnitPrice    float64   `json:"unit_price"`
	Subtotal     float64   `json:"subtotal"`
}

type TimelineResponse struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	User        string    `json:"user"`
	Timestamp   time.Time `json:"timestamp"`
}

func (o *Order) ToResponse() OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemResponse{
			ID:           item.ID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductSKU:   item.ProductSKU,
			ProductImage: item.ProductImage,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			Subtotal:     item.Subtotal,
		}
	}

	timeline := make([]TimelineResponse, len(o.Timeline))
	for i, t := range o.Timeline {
		timeline[i] = TimelineResponse{
			ID:          t.ID,
			Type:        t.Type,
			Title:       t.Title,
			Description: t.Description,
			User:        t.UserName,
			Timestamp:   t.CreatedAt,
		}
	}

	customer := CustomerBrief{Name: o.CustomerName}
	if o.CustomerID != nil {
		customer.ID = o.CustomerID.String()
		customer.Email = o.CustomerEmail
	}

	return OrderResponse{
		ID:              o.ID,
		OrderNumber:     o.OrderNumber,
		Customer:        customer,
		Status:          o.Status,
		PaymentStatus:   o.PaymentStatus,
		PaymentMethodID: o.PaymentMethodID,
		PaymentMethod:   o.PaymentMethodName,
		ShippingMethod:  o.ShippingMethod,
		ShippingAddress: o.ShippingAddress,
		Subtotal:        o.Subtotal,
		ShippingCost:    o.ShippingCost,
		Tax:             o.Tax,
		Discount:        o.Discount,
		TotalAmount:     o.TotalAmount,
		Notes:           o.Notes,
		Items:           items,
		Timeline:        timeline,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}
