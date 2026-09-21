package category

import (
	"time"

	"github.com/google/uuid"
)

// Category represents the `categories` table.
type Category struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  uuid.UUID  `gorm:"column:tenant_id;type:uuid;not null;index"`
	ParentID  *uuid.UUID `gorm:"column:parent_id;type:uuid"`
	Name      string     `gorm:"column:name;not null"`
	Slug      string     `gorm:"column:slug;not null"`
	Description string   `gorm:"column:description"`
	Image     string     `gorm:"column:image"`
	IsActive  bool       `gorm:"column:is_active;default:true"`
	SortOrder int        `gorm:"column:sort_order;default:0"`
	CreatedBy string     `gorm:"column:created_by"`
	UpdatedBy string     `gorm:"column:updated_by"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Category) TableName() string {
	return "categories"
}

type CreateCategoryRequest struct {
	ParentID  *uuid.UUID `json:"parent_id"`
	Name      string     `json:"name" binding:"required,min=2,max=200"`
	Slug      string     `json:"slug" binding:"omitempty,min=2,max=200"`
	Description string   `json:"description"`
	Image     string     `json:"image"`
	SortOrder int        `json:"sort_order"`
	IsActive  *bool      `json:"is_active"`
}

type UpdateCategoryRequest struct {
	ParentID  *uuid.UUID `json:"parent_id"`
	Name      *string    `json:"name" binding:"omitempty,min=2,max=200"`
	Slug      *string    `json:"slug" binding:"omitempty,min=2,max=200"`
	Description *string  `json:"description"`
	Image     *string    `json:"image"`
	SortOrder *int       `json:"sort_order"`
	IsActive  *bool      `json:"is_active"`
}

type CategoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Description string   `json:"description"`
	Image     string     `json:"image"`
	IsActive  bool       `json:"is_active"`
	SortOrder int        `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		ParentID:    c.ParentID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Image:       c.Image,
		IsActive:    c.IsActive,
		SortOrder:   c.SortOrder,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
