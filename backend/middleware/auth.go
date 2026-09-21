package middleware

import (
	"backend-go/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Auth validates the JWT access token from the Authorization header.
// It sets "UserId", "UserGroupId", "TenantId", "Role", "IsAdministrator",
// and "IsPlatform" in the gin context.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := token.VerifyToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized: " + err.Error(),
			})
			c.Abort()
			return
		}

		userID, _ := claims["UserId"].(string)
		userGroupID, _ := claims["UserGroupId"].(string)
		tenantID, _ := claims["TenantId"].(string)
		role, _ := claims["Role"].(string)
		isAdmin, _ := claims["IsAdministrator"].(bool)
		application, _ := claims["Application"].(string)

		c.Set("UserId", userID)
		c.Set("UserGroupId", userGroupID)
		c.Set("TenantId", tenantID)
		c.Set("Role", role)
		c.Set("IsAdministrator", isAdmin)
		c.Set("IsPlatform", tenantID == "" || tenantID == "00000000-0000-0000-0000-000000000000")
		c.Set("Application", application)

		c.Next()
	}
}

// APIKey validates the X-API-KEY header for service-to-service calls.
func APIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := token.VerifyAPIKey(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized: " + err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("AuthType", "api_key")
		c.Next()
	}
}
