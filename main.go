package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goldcast/gc_auth_service/internal/config"
	"github.com/goldcast/gc_auth_service/internal/handlers"
	"github.com/goldcast/gc_auth_service/internal/middleware"
	"github.com/goldcast/gc_auth_service/internal/routes"
	"github.com/goldcast/gc_auth_service/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(logger)

	// Setup routes
	routes.SetupRoutes(router, authHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Infof("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
