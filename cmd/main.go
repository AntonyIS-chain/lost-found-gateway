package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/middlewares"
	"github.com/AntonyIS-chain/lost-found-gateway/pkg"
	"github.com/gin-gonic/gin"
)

// GatewayServer initializes and starts the API Gateway
func GatewayServer() {
	cfg := config.LoadConfig() // Load configuration

	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	// Apply global middleware
	router.Use(middlewares.CORSMiddleware())
	// router.Use(middlewares.LoggingMiddleware())
	// router.Use(middlewares.RateLimiterMiddleware())

	// Public Routes (No Authentication)
	router.POST("/api/v1/auth/login", pkg.NewReverseProxy(cfg.USER_SERVICE))
	router.POST("/api/v1/auth/register", pkg.NewReverseProxy(cfg.USER_SERVICE))
	router.POST("/api/v1/auth/refresh-token", pkg.NewReverseProxy(cfg.USER_SERVICE))
	router.POST("/api/v1/auth/reset-password", pkg.NewReverseProxy(cfg.USER_SERVICE))

	// Protected Routes (Require valid JWT)
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middlewares.JWTMiddleware())

	// User Service Routes
	userRoutes := protectedRoutes.Group("/api/v1/users")
	adapters.RegisterProxyRoutes(userRoutes, cfg.USER_SERVICE)

	// Matching Service Routes
	matchingRoutes := protectedRoutes.Group("/api/v1/matching")
	adapters.RegisterProxyRoutes(matchingRoutes, cfg.MATCHING_SERVICE)

	// Reward Service Routes
	rewardRoutes := protectedRoutes.Group("/reward")
	adapters.RegisterProxyRoutes(rewardRoutes, cfg.REWARD_SERVICE)

	// Payment Service Routes
	paymentRoutes := protectedRoutes.Group("/payment")
	adapters.RegisterProxyRoutes(paymentRoutes, cfg.PAYMENT_SERVICE)

	// Notification Service Routes
	notificationRoutes := protectedRoutes.Group("/notification")
	adapters.RegisterProxyRoutes(notificationRoutes, cfg.NOTIFICATION_SERVICE)

	// Document Service Routes
	documentRoutes := protectedRoutes.Group("/document")
	adapters.RegisterProxyRoutes(documentRoutes, cfg.DOCUMENT_SERVICE)

	// Start the HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.GATEWAY_PORT),
		Handler: router,
	}

	// Graceful Shutdown
	go func() {
		log.Printf("[INFO] API Gateway running on port %s\n", cfg.GATEWAY_PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("[INFO] Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server shutdown failed: %v\n", err)
	}

	log.Println("[INFO] API Gateway stopped successfully.")
}
