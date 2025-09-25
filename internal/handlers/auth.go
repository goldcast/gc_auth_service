package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goldcast/gc_auth_service/internal/models"
	"github.com/goldcast/gc_auth_service/pkg/jwt"
	"github.com/goldcast/gc_auth_service/pkg/logger"
	"github.com/goldcast/gc_auth_service/pkg/password"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	logger     *logger.Logger
	jwtService *jwt.Service
	validator  *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		logger:     logger,
		jwtService: jwt.New("your-secret-key-change-in-production", 24), // TODO: Get from config
		validator:  validator.New(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithField("error", err.Error()).Warn("Invalid registration request")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to hash password")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	// Create user (in a real app, this would be saved to database)
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
	}).Info("User registered successfully")

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithField("error", err.Error()).Warn("Invalid login request")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// TODO: In a real app, fetch user from database
	// For now, we'll simulate a user
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  "testuser",
		Password:  "$2a$10$example_hash", // This would be the actual hashed password from DB
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check password (in real app, compare with stored hash)
	// For demo purposes, we'll accept any password
	// if !password.CheckPasswordHash(req.Password, user.Password) {
	//     c.JSON(http.StatusUnauthorized, models.APIResponse{
	//         Success: false,
	//         Message: "Invalid credentials",
	//     })
	//     return
	// }

	// Generate tokens
	accessToken, err := h.jwtService.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to generate access token")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to generate refresh token")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	// Remove password from response
	user.Password = ""

	response := models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in successfully")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithField("error", err.Error()).Warn("Invalid refresh token request")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		h.logger.WithField("error", err.Error()).Warn("Invalid refresh token")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid or expired refresh token",
		})
		return
	}

	// TODO: In a real app, verify the refresh token exists in database
	// and fetch user details

	// Generate new access token
	accessToken, err := h.jwtService.GenerateToken(claims.UserID, claims.Email, claims.Username)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to generate new access token")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	response := map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   24 * 60 * 60, // 24 hours in seconds
	}

	h.logger.WithField("user_id", claims.UserID).Info("Token refreshed successfully")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// TODO: In a real app, fetch user from database
	user := &models.User{
		ID:        userID.(uuid.UUID),
		Email:     c.GetString("user_email"),
		Username:  c.GetString("user_username"),
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: In a real app, invalidate the refresh token in database
	// For now, we'll just return success

	h.logger.WithField("user_id", c.GetString("user_id")).Info("User logged out")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}
