package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goldcast/gc_auth_service/internal/models"
	"github.com/goldcast/gc_auth_service/pkg/jwt"
	"github.com/goldcast/gc_auth_service/pkg/logger"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(log *logger.Logger, jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			log.WithField("error", err.Error()).Warn("Invalid token")
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		c.Next()
	}
}
