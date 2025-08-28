package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole membatasi akses berdasarkan role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found in context"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowed := range roles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
		c.Abort()
	}
}
