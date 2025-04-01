package app

import (
	"fmt"
	"log"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/app/controllers"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/adapters/middlewares"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/AntonyIS-chain/lost-found-gateway/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterProxyRoutes sets up reverse proxy routes for a given service
func RegisterProxyRoutes(router *gin.RouterGroup, target string) {
	router.Any("/*proxyPath", func(c *gin.Context) {
		proxyPath := c.Param("proxyPath")
		if proxyPath == "" {
			proxyPath = "/"
		}

		pkg.NewReverseProxy(target)(c)
	})
}

func InitGinRoutes(userSvc ports.AuthService, config *config.Config) {
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
	authRoutes := router.Group("/v1/api/auth")
	{
		authRoutes.POST("/login", userController.Authenticate)
		authRoutes.POST("/refresh-token", userController.RefreshToken)
		authRoutes.POST("/validate-token", userController.ValidateToken)
		authRoutes.POST("/logout", userController.InValidateToken)
	}

	// Protected Routes (Require valid JWT)
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middlewares.JWTMiddleware())

	// User Service Routes
	userRoutes := protectedRoutes.Group("/v1/api/users")
	RegisterProxyRoutes(userRoutes, config.USER_SERVICE)

	// Matching Service Routes
	matchingRoutes := protectedRoutes.Group("/v1/api/matching")
	RegisterProxyRoutes(matchingRoutes, config.MATCHING_SERVICE)

	// Reward Service Routes
	rewardRoutes := protectedRoutes.Group("/reward")
	RegisterProxyRoutes(rewardRoutes, config.REWARD_SERVICE)

	// Payment Service Routes
	paymentRoutes := protectedRoutes.Group("/payment")
	RegisterProxyRoutes(paymentRoutes, config.PAYMENT_SERVICE)

	// Notification Service Routes
	notificationRoutes := protectedRoutes.Group("/notification")
	RegisterProxyRoutes(notificationRoutes, config.NOTIFICATION_SERVICE)

	// Document Service Routes
	documentRoutes := protectedRoutes.Group("/document")
	RegisterProxyRoutes(documentRoutes, config.DOCUMENT_SERVICE)
	// Start server
	log.Println("Starting server on port", config.GATEWAY_PORT)
	router.Run(fmt.Sprintf(":%s", config.GATEWAY_PORT))
}
