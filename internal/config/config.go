package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Environment string
	Port        string
	LogLevel    string
	JWTSecret   string
	JWTExpiry   int // in hours
	DatabaseURL string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:   getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
