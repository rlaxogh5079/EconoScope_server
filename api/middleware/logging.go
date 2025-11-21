package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rlaxogh5079/EconoScope/pkg/logger"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Log.WithFields(map[string]interface{}{
			"status":   status,
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"ip":       c.ClientIP(),
			"latency":  latency.String(),
			"userAgent": c.Request.UserAgent(),
		}).Info("HTTP Request")
	}
}