package product

import (
	"backend-go/internal/entities"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Database Models
// =============================================================================

// Product represents the `products` table.
type Product struct {
	ID                uuid.UUID             `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID          uuid.UUID             `gorm:"column:tenant_id;type:uuid;not null;index"`
	Code              string                `gorm:"column:code"`
	SKU               string                `gorm:"column:sku"`
	Name              string                `gorm:"column:name;not null"`
	Description       string                `gorm:"column:description"`
	Price             float64               `gorm:"column:price;not null;default:0"`
	CompareAtPrice    float64               `gorm:"column:compare_at_price;not null;default:0"`
	CostPrice         float64               `gorm:"column:cost_price;not null;default:0"`
	Stock             int                   `gorm:"column:stock;not null;default:0"`
	Reserved          int                   `gorm:"column:reserved;not null;default:0"`
	Incoming          int                   `gorm:"column:incoming;not null;default:0"`
	LowStockThreshold int                   `gorm:"column:low_stock_threshold;not null;default:5"`
	Status            string                `gorm:"column:status;not null;default:'active'"`
	CategoryID        *uuid.UUID            `gorm:"column:category_id;type:uuid"`
	Tags              entities.StringList   `gorm:"column:tags;type:jsonb"`
	IsActive          bool                  `gorm:"column:is_active;default:true"`
	Rating            float64               `gorm:"column:rating;not null;default:0"`
	ReviewCount       int                   `gorm:"column:review_count;not null;default:0"`
	Images            []ProductImage        `gorm:"foreignKey:ProductID"`
	CreatedBy         string                `gorm:"column:created_by"`
	UpdatedBy         string                `gorm:"column:updated_by"`
	CreatedAt         time.Time             `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time             `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt         *time.Time            `gorm:"column:deleted_at"`

	// Transient field (joined, not stored)
	CategoryName string `gorm:"-"`
}

func (Product) TableName() string {
	return "products"
}

// ProductImage represents the `product_images` table.
type ProductImage struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID uuid.UUID `gorm:"column:product_id;type:uuid;not null;index"`
	URL       string    `gorm:"column:url;not null"`
	Alt       string    `gorm:"column:alt"`
	SortOrder int       `gorm:"column:sort_order;not null;default:0"`
	IsPrimary bool      `gorm:"column:is_primary;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

// Product statuses
const (
	StatusActive   = "active"
	StatusDraft    = "draft"
	StatusArchived = "archived"
)

// =============================================================================
// Request DTOs
// =============================================================================

type CreateProductRequest struct {
	SKU               string     `json:"sku" binding:"required,min=2,max=100"`
	Name              string     `json:"name" binding:"required,min=3,max=200"`
	Description       string     `json:"description"`
	Price             float64    `json:"price" binding:"required,gt=0"`
	CompareAtPrice    float64    `json:"compare_at_price"`
	CostPrice         float64    `json:"cost_price"`
	Stock             int        `json:"stock" binding:"gte=0"`
	LowStockThreshold int        `json:"low_stock_threshold"`
	Status            string     `json:"status" binding:"omitempty,oneof=active draft archived"`
	CategoryID        *uuid.UUID `json:"category_id"`
	Tags              []string   `json:"tags"`
	Images            []string   `json:"images"`
}

type UpdateProductRequest struct {
	SKU               *string    `json:"sku" binding:"omitempty,min=2,max=100"`
	Name              *string    `json:"name" binding:"omitempty,min=3,max=200"`
	Description       *string    `json:"description"`
	Price             *float64   `json:"price" binding:"omitempty,gt=0"`
	CompareAtPrice    *float64   `json:"compare_at_price"`
	CostPrice         *float64   `json:"cost_price"`
	LowStockThreshold *int       `json:"low_stock_threshold"`
	Status            *string    `json:"status" binding:"omitempty,oneof=active draft archived"`
	CategoryID        *uuid.UUID `json:"category_id"`
	Tags              []string   `json:"tags"`
	Images            []string   `json:"images"`
	IsActive          *bool      `json:"is_active"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type ProductResponse struct {
	ID                uuid.UUID       `json:"id"`
	SKU               string          `json:"sku"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Price             float64         `json:"price"`
	CompareAtPrice    float64         `json:"compare_at_price"`
	CostPrice         float64         `json:"cost_price"`
	Stock             int             `json:"stock"`
	Reserved          int             `json:"reserved"`
	Incoming          int             `json:"incoming"`
	LowStockThreshold int             `json:"low_stock_threshold"`
	Status            string          `json:"status"`
	CategoryID        *uuid.UUID      `json:"category_id,omitempty"`
	CategoryName      string          `json:"category_name"`
	Tags              []string        `json:"tags"`
	IsActive          bool            `json:"is_active"`
	Rating            float64         `json:"rating"`
	ReviewCount       int             `json:"review_count"`
	Images            []ImageResponse `json:"images"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type ImageResponse struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Alt       string    `json:"alt"`
	IsPrimary bool      `json:"is_primary"`
}

func (p *Product) ToResponse() ProductResponse {
	images := make([]ImageResponse, len(p.Images))
	for i, img := range p.Images {
		images[i] = ImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			Alt:       img.Alt,
			IsPrimary: img.IsPrimary,
		}
	}

	return ProductResponse{
		ID:                p.ID,
		SKU:               p.SKU,
		Name:              p.Name,
		Description:       p.Description,
		Price:             p.Price,
		CompareAtPrice:    p.CompareAtPrice,
		CostPrice:         p.CostPrice,
		Stock:             p.Stock,
		Reserved:          p.Reserved,
		Incoming:          p.Incoming,
		LowStockThreshold: p.LowStockThreshold,
		Status:            p.Status,
		CategoryID:        p.CategoryID,
		CategoryName:      p.CategoryName,
		Tags:              p.Tags,
		IsActive:          p.IsActive,
		Rating:            p.Rating,
		ReviewCount:       p.ReviewCount,
		Images:            images,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}
