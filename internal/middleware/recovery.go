package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goldcast/gc_auth_service/pkg/logger"
)

// Recovery returns a gin.HandlerFunc for recovering from panics
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.WithField("error", err).Error("Panic recovered")
		}
		
		// Check for a broken connection
		if ne, ok := recovered.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
					strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					c.Abort()
					return
				}
			}
		}

		// Log the stack trace
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		log.WithFields(map[string]interface{}{
			"error":   recovered,
			"request": string(httpRequest),
			"stack":   string(debug.Stack()),
		}).Error("Panic recovered")

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
