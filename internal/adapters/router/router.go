package router

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Auth routes (proxied to Auth Service)
	authServiceURL := "http://auth-service:8081" // Replace with your real service URL or config

	router.POST("/auth/signup", func(c *gin.Context) {
		proxyRequest(c, authServiceURL+"/signup")
	})

	router.POST("/auth/login", func(c *gin.Context) {
		proxyRequest(c, authServiceURL+"/login")
	})

	router.POST("/auth/forgot", func(c *gin.Context) {
		proxyRequest(c, authServiceURL+"/forgot")
	})

	router.POST("/auth/reset", func(c *gin.Context) {
		proxyRequest(c, authServiceURL+"/reset")
	})
}

// 🔁 Shared proxy logic
func proxyRequest(c *gin.Context, targetURL string) {
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "auth service unreachable"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}
