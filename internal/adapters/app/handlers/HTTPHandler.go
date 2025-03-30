package app

import (
	"fmt"
	"log"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/app/controllers"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/middlewares"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(userSvc ports.AuthenticationService, config *config.Config) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()


	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middlewares.CORSMiddleware())
	// router.Use(middlewares.LoggingMiddleware())
	// router.Use(middlewares.RateLimiterMiddleware())

	// Initialize Controllers
	userController := controllers.NewAuthenticationController(userSvc)

	// User Routes
	authRoutes := router.Group("/v1/gw/auth")
	{
		authRoutes.POST("/login", userController.Authenticate)
		authRoutes.POST("/signup", userController.Authenticate)
		authRoutes.POST("/refresh-token", userController.RefreshToken)
		authRoutes.POST("/validate-token", userController.ValidateToken)
		authRoutes.POST("/reset-password", userController.ValidateToken)
	}

	// Protected Routes (Require valid JWT)
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middlewares.JWTMiddleware())

	// User Service Routes
	userRoutes := protectedRoutes.Group("/v1/api/users")
	adapters.RegisterProxyRoutes(userRoutes, config.USER_SERVICE)

	// Matching Service Routes
	matchingRoutes := protectedRoutes.Group("/v1/api/matching")
	adapters.RegisterProxyRoutes(matchingRoutes, config.MATCHING_SERVICE)

	// Reward Service Routes
	rewardRoutes := protectedRoutes.Group("/reward")
	adapters.RegisterProxyRoutes(rewardRoutes, config.REWARD_SERVICE)

	// Payment Service Routes
	paymentRoutes := protectedRoutes.Group("/payment")
	adapters.RegisterProxyRoutes(paymentRoutes, config.PAYMENT_SERVICE)

	// Notification Service Routes
	notificationRoutes := protectedRoutes.Group("/notification")
	adapters.RegisterProxyRoutes(notificationRoutes, config.NOTIFICATION_SERVICE)

	// Document Service Routes
	documentRoutes := protectedRoutes.Group("/document")
	adapters.RegisterProxyRoutes(documentRoutes, config.DOCUMENT_SERVICE)
	// Start server
	log.Println("Starting server on port", config.GATEWAY_PORT)
	router.Run(fmt.Sprintf(":%s", config.GATEWAY_PORT))
}
