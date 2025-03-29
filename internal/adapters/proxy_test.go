package adapters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Custom ResponseWriter to avoid CloseNotify issue
type customResponseWriter struct {
	*httptest.ResponseRecorder
}

func (c *customResponseWriter) CloseNotify() <-chan bool {
	// Return a dummy channel to prevent issues with CloseNotify
	ch := make(chan bool, 1)
	return ch
}

// Mock target service handler
func mockTargetService() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock service response"))
	})
	return httptest.NewServer(handler)
}

func TestRegisterProxyRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Start a mock backend service
	mockServer := mockTargetService()
	defer mockServer.Close()

	// Initialize router and register proxy
	router := gin.Default()
	apiGroup := router.Group("/api")
	RegisterProxyRoutes(apiGroup, mockServer.URL)

	t.Run("Should forward request to backend service", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		w := &customResponseWriter{httptest.NewRecorder()} // Use custom ResponseWriter

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Mock service response")
	})

	t.Run("Should handle non-existent routes gracefully", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/non-existent", nil)
		w := &customResponseWriter{httptest.NewRecorder()} // Use custom ResponseWriter

		router.ServeHTTP(w, req)

		// Expected 404 if backend doesn't handle it
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
