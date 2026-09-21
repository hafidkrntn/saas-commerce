package review

import (
	"backend-go/pkg/pagination"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review represents the `reviews` table.
type Review struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID       uuid.UUID  `gorm:"column:tenant_id;type:uuid;not null;index"`
	ProductID      uuid.UUID  `gorm:"column:product_id;type:uuid;not null;index"`
	CustomerID     *uuid.UUID `gorm:"column:customer_id;type:uuid"`
	CustomerName   string     `gorm:"column:customer_name;not null"`
	CustomerAvatar string     `gorm:"column:customer_avatar"`
	Rating         int        `gorm:"column:rating;not null"`
	Title          string     `gorm:"column:title"`
	Content        string     `gorm:"column:content"`
	IsVerified     bool       `gorm:"column:is_verified;not null;default:false"`
	IsPublished    bool       `gorm:"column:is_published;not null;default:false"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Review) TableName() string {
	return "reviews"
}

type CreateReviewRequest struct {
	ProductID  uuid.UUID `json:"product_id" binding:"required"`
	CustomerID *uuid.UUID `json:"customer_id"`
	CustomerName string  `json:"customer_name" binding:"required,min=2,max=200"`
	Rating     int       `json:"rating" binding:"required,min=1,max=5"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
}

type ModerateReviewRequest struct {
	IsPublished *bool `json:"is_published"`
	IsVerified  *bool `json:"is_verified"`
}

type ReviewResponse struct {
	ID             uuid.UUID  `json:"id"`
	ProductID      uuid.UUID  `json:"product_id"`
	CustomerID     *uuid.UUID `json:"customer_id,omitempty"`
	CustomerName   string     `json:"customer_name"`
	CustomerAvatar string     `json:"customer_avatar"`
	Rating         int        `json:"rating"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	IsVerified     bool       `json:"is_verified"`
	IsPublished    bool       `json:"is_published"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (r *Review) ToResponse() ReviewResponse {
	return ReviewResponse{
		ID:             r.ID,
		ProductID:      r.ProductID,
		CustomerID:     r.CustomerID,
		CustomerName:   r.CustomerName,
		CustomerAvatar: r.CustomerAvatar,
		Rating:         r.Rating,
		Title:          r.Title,
		Content:        r.Content,
		IsVerified:     r.IsVerified,
		IsPublished:    r.IsPublished,
		CreatedAt:      r.CreatedAt,
	}
}

// RepositoryInterface defines data access for reviews.
type RepositoryInterface interface {
	Create(tx *gorm.DB, r *Review) error
	GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Review, error)
	List(tx *gorm.DB, tenantID uuid.UUID, productID string, onlyPublished bool, page, limit int) (pagination.Response[Review], error)
	Update(tx *gorm.DB, r *Review) error
	Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	// AggregateRating returns the average rating and count for a product.
	AggregateRating(tx *gorm.DB, tenantID uuid.UUID, productID uuid.UUID) (float64, int64, error)
}
