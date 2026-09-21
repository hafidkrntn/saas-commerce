package tenant

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents the `tenants` table.
type Tenant struct {
	ID                 uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name               string     `gorm:"column:name;not null"`
	Slug               string     `gorm:"column:slug;not null"`
	Logo               string     `gorm:"column:logo"`
	PlanID             *uuid.UUID `gorm:"column:plan_id;type:uuid"`
	SubscriptionStatus string     `gorm:"column:subscription_status;not null;default:'trial'"`
	TrialEndsAt        *time.Time `gorm:"column:trial_ends_at"`
	CreatedBy          string     `gorm:"column:created_by"`
	UpdatedBy          string     `gorm:"column:updated_by"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// Plan represents the `plans` table.
type Plan struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name"`
	Code        string    `gorm:"column:code"`
	Price       float64   `gorm:"column:price"`
	Description string    `gorm:"column:description"`
}

func (Plan) TableName() string {
	return "plans"
}

// =============================================================================
// Request DTOs
// =============================================================================

type CreateTenantRequest struct {
	Name               string     `json:"name" binding:"required,min=2,max=200"`
	Slug               string     `json:"slug" binding:"required,min=2,max=100"`
	Logo               string     `json:"logo"`
	PlanID             *uuid.UUID `json:"plan_id"`
	SubscriptionStatus string     `json:"subscription_status" binding:"omitempty,oneof=trial active past_due cancelled"`
	TrialEndsAt        *time.Time `json:"trial_ends_at"`
}

type UpdateTenantRequest struct {
	Name               *string    `json:"name" binding:"omitempty,min=2,max=200"`
	Slug               *string    `json:"slug" binding:"omitempty,min=2,max=100"`
	Logo               *string    `json:"logo"`
	PlanID             *uuid.UUID `json:"plan_id"`
	SubscriptionStatus *string    `json:"subscription_status" binding:"omitempty,oneof=trial active past_due cancelled"`
	TrialEndsAt        *time.Time `json:"trial_ends_at"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type TenantResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Slug               string     `json:"slug"`
	Logo               string     `json:"logo"`
	PlanID             *uuid.UUID `json:"plan_id,omitempty"`
	PlanCode           string     `json:"plan_code"`
	SubscriptionStatus string     `json:"subscription_status"`
	TrialEndsAt        *time.Time `json:"trial_ends_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

type PlanResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
}
