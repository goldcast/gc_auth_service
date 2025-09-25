package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goldcast/gc_auth_service/pkg/logger"
)

// Logger returns a gin.HandlerFunc for logging requests
func Logger(log *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.WithFields(map[string]interface{}{
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
			"status":      param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"user_agent":  param.Request.UserAgent(),
			"error":       param.ErrorMessage,
		}).Info("HTTP Request")
		return ""
	})
}
