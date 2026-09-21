package middleware

import (
	"backend-go/pkg/apperror"
	"backend-go/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Recovery catches panics and returns a 500 JSON response instead of crashing.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		msg := fmt.Sprintf("panic recovered: %v", recovered)
		err := apperror.Internal(msg, nil)
		response.Err(c, err)
	})
}
