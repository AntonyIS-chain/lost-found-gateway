package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test ValidateBasicAuth function
func TestValidateBasicAuth(t *testing.T) {
	tests := []struct {
		username string
		password string
		expected bool
	}{
		{"admin", "password123", true}, // ✅ Correct credentials
		{"user1", "mypassword", true},  // ✅ Correct credentials
		{"admin", "wrongpass", false},  // ❌ Wrong password
		{"unknown", "password123", false}, // ❌ User doesn't exist
	}

	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			result := ValidateBasicAuth(test.username, test.password)
			assert.Equal(t, test.expected, result)
		})
	}
}

// Test AuthorizationMiddleware
func TestAuthorizationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		username     string
		expectedCode int
		requiredRole string
	}{
		{"admin", http.StatusOK, "admin"},     // ✅ Admin has access
		{"user1", http.StatusOK, "user"},      // ✅ User has access
		{"user1", http.StatusForbidden, "admin"}, // ❌ User trying admin access
		{"unknown", http.StatusUnauthorized, "admin"}, // ❌ Unknown user
	}

	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				if test.username != "unknown" {
					c.Set("username", test.username) // Simulate user login
				}
				c.Next()
			})
			router.GET("/protected", AuthorizationMiddleware(test.requiredRole), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Success"})
			})

			req, _ := http.NewRequest("GET", "/protected", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}
