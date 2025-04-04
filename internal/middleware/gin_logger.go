package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger is a middleware that logs requests using a Zap logger.
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process the request
		c.Next()

		// Capture the status and latency after the request is processed
		status := c.Writer.Status()
		latency := time.Since(start)

		// Log the HTTP request and response details
		logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		)
	}
}

// GinRecovery logs any panics, preventing the server from crashing.
func GinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic
		logger.Error("HTTP panic recovered",
			zap.Any("error", recovered),
		)

		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		c.Abort()
	})
}
