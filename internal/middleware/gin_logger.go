package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger returns a Gin middleware handler that logs HTTP requests and responses using the provided Zap logger.
//
// The middleware captures key details of incoming HTTP requests, including:
//   - Request method (e.g., GET, POST)
//   - Request path (endpoint)
//   - Response HTTP status code
//   - Request processing latency
//
// These structured logs greatly assist developers and operators with monitoring, debugging, and analysis of request patterns and performance characteristics.
//
// Parameters:
//   - logger: A Zap logger instance used to emit structured logs.
//
// Returns:
//   - gin.HandlerFunc: Middleware handler function suitable for inclusion in a Gin router's middleware chain.
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Proceed with request processing
		c.Next()

		// Record status code and latency after request completion
		status := c.Writer.Status()
		latency := time.Since(start)

		// Log structured request and response details
		logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		)
	}
}

// GinRecovery returns a Gin middleware handler that recovers from panics, preventing server crashes, and logs any recovered panics using the provided Zap logger.
//
// In the event of a panic during request processing, this middleware:
//   - Captures and recovers from the panic gracefully.
//   - Logs the panic details as structured error messages.
//   - Returns a standard HTTP 500 (Internal Server Error) JSON response to the client.
//
// This approach ensures higher reliability of the HTTP server and simplifies troubleshooting by providing detailed, structured logging around panic events.
//
// Parameters:
//   - logger: A Zap logger instance used to log recovered panics.
//
// Returns:
//   - gin.HandlerFunc: Middleware handler function suitable for inclusion in a Gin router's middleware chain.
func GinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log detailed information about the recovered panic
		logger.Error("HTTP panic recovered",
			zap.Any("error", recovered),
		)

		// Respond with a standard HTTP 500 error message
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})

		c.Abort()
	})
}
