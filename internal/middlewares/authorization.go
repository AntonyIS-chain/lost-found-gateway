package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var validUsers = map[string]string{
	"admin": "password123",
	"user1": "mypassword",
}

// ✅ ValidateBasicAuth checks username and password
func ValidateBasicAuth(username, password string) bool {
	if pass, exists := validUsers[username]; exists && pass == password {
		return true
	}
	return false
}

// ✅ AuthorizationMiddleware validates JWT token
func AuthorizationMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Define roles (you can replace this with a DB lookup)
		userRoles := map[string]string{
			"admin": "admin",
			"user1": "user",
		}

		role, ok := userRoles[username.(string)]
		if !ok || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}
