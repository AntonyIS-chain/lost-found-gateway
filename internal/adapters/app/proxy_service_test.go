package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestForwardRequest(t *testing.T) {
	// Mock backend service
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/mock-path", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Mock response"}`))
	}))
	defer mockServer.Close()

	// Create a ProxyService instance
	proxyService := NewProxyService()

	// Prepare the request
	proxyReq := domain.ProxyRequest{
		Method: "GET",
		Path:   mockServer.URL + "/mock-path", // Use mock server URL
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Call the ForwardRequest method
	response, err := proxyService.ForwardRequest(proxyReq)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Headers["Content-Type"])
	assert.JSONEq(t, `{"message":"Mock response"}`, string(response.Body))
}
