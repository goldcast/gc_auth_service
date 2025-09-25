package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goldcast/gc_auth_service/internal/handlers"
	"github.com/goldcast/gc_auth_service/internal/middleware"
	"github.com/goldcast/gc_auth_service/pkg/jwt"
	"github.com/goldcast/gc_auth_service/pkg/logger"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "gc_auth_service",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		{
			// Initialize JWT service for middleware
			jwtService := jwt.New("your-secret-key-change-in-production", 24) // TODO: Get from config
			logger := logger.New("info") // TODO: Get from config
			
			protected.Use(middleware.AuthMiddleware(logger, jwtService))
			{
				protected.GET("/profile", authHandler.GetProfile)
				protected.POST("/logout", authHandler.Logout)
			}
		}
	}
}
