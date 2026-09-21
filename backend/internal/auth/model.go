package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserRow mirrors the `users` table for the auth flow.
type UserRow struct {
	ID           uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     *uuid.UUID `gorm:"column:tenant_id;type:uuid"`
	Name         string     `gorm:"column:name"`
	Email        string     `gorm:"column:email"`
	PasswordHash string     `gorm:"column:password_hash"`
	Avatar       string     `gorm:"column:avatar"`
	Status       string     `gorm:"column:status"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
}

func (UserRow) TableName() string {
	return "users"
}

// RoleRow mirrors the `roles` table for the auth flow.
type RoleRow struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID *uuid.UUID `gorm:"column:tenant_id;type:uuid"`
	Name     string    `gorm:"column:name"`
	Permissions int64  `gorm:"column:permissions"`
}

func (RoleRow) TableName() string {
	return "roles"
}

// RefreshToken mirrors the `refresh_tokens` table.
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `gorm:"column:user_id;type:uuid"`
	TokenHash string     `gorm:"column:token_hash"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// =============================================================================
// Request DTOs
// =============================================================================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	Role      string     `json:"role"`
	IsAdmin   bool       `json:"is_admin"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type TenantBrief struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	PlanCode    string    `json:"plan_code"`
	SubscriptionStatus string `json:"subscription_status"`
}
