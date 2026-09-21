package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents the `users` table (tenant-scoped).
type User struct {
	ID           uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     *uuid.UUID `gorm:"column:tenant_id;type:uuid;index"`
	Name         string     `gorm:"column:name;not null"`
	Email        string     `gorm:"column:email;not null"`
	PasswordHash string     `gorm:"column:password_hash;not null"`
	Avatar       string     `gorm:"column:avatar"`
	Status       string     `gorm:"column:status;not null;default:'active'"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedBy    string     `gorm:"column:created_by"`
	UpdatedBy    string     `gorm:"column:updated_by"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// Role represents the `roles` table.
type Role struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    *uuid.UUID `gorm:"column:tenant_id;type:uuid;index"`
	Name        string     `gorm:"column:name;not null"`
	Description string     `gorm:"column:description"`
	Permissions int64      `gorm:"column:permissions;not null;default:0"`
	IsSystem    bool       `gorm:"column:is_system;not null;default:false"`
	CreatedBy   string     `gorm:"column:created_by"`
	UpdatedBy   string     `gorm:"column:updated_by"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Role) TableName() string {
	return "roles"
}

// UserRole joins users to roles.
type UserRole struct {
	UserID uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	RoleID uuid.UUID `gorm:"column:role_id;type:uuid;primaryKey"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// User statuses
const (
	StatusActive  = "active"
	StatusInactive = "inactive"
	StatusInvited = "invited"
)

// =============================================================================
// Request DTOs
// =============================================================================

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=200"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id" binding:"required"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=2,max=200"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive invited"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Permissions int64  `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description"`
	Permissions *int64  `json:"permissions"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

type RoleResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Permissions int64      `json:"permissions"`
	Users       int64      `json:"users"`
	IsSystem    bool       `json:"is_system"`
	CreatedAt   time.Time  `json:"created_at"`
}
