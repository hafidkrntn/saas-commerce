package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole rejects requests whose token role is not in the allowed list.
// An empty list allows any authenticated user.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		if len(allowed) == 0 {
			c.Next()
			return
		}

		role := c.GetString("Role")
		if isAdmin := c.GetBool("IsAdministrator"); isAdmin && allowed["admin"] {
			c.Next()
			return
		}

		if !allowed[role] {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "akses ditolak: role tidak memiliki izin",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePlatform restricts routes to platform-admin tokens (no tenant scope).
func RequirePlatform() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("IsPlatform") {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "akses ditolak: khusus platform admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireTenant restricts routes to tenant-scoped tokens.
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("TenantId") == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "akses ditolak: tenant scope diperlukan",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
