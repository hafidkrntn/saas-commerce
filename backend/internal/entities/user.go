package entities

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User represents the authenticated user information extracted from the token.
type User struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TenantID        string `json:"tenant_id"`
	Role            string `json:"role"`
	IsAdministrator bool   `json:"is_administrator"`
}

// TenantContext carries the resolved tenant/user scope for a request.
type TenantContext struct {
	TenantID    uuid.UUID
	UserID      string
	Role        string
	IsPlatform  bool
	IsAdmin     bool
	Application string
}

// GetTenantContext reads the tenant/user scope set by middleware.Auth.
func GetTenantContext(c *gin.Context) TenantContext {
	tenantID, err := uuid.Parse(c.GetString("TenantId"))
	if err != nil {
		tenantID = uuid.Nil
	}

	return TenantContext{
		TenantID:    tenantID,
		UserID:      c.GetString("UserId"),
		Role:        c.GetString("Role"),
		IsPlatform:  c.GetBool("IsPlatform"),
		IsAdmin:     c.GetBool("IsAdministrator"),
		Application: c.GetString("Application"),
	}
}

// TenantID returns the tenant id from the gin context (nil if not scoped).
func TenantID(c *gin.Context) uuid.UUID {
	return GetTenantContext(c).TenantID
}

// RequireTenantID returns an error when the request has no tenant scope
// (i.e. a platform-admin token hitting a tenant-scoped endpoint).
func (t TenantContext) RequireTenantID() (uuid.UUID, error) {
	if t.TenantID == uuid.Nil {
		return uuid.Nil, errors.New("tenant scope required")
	}
	return t.TenantID, nil
}
