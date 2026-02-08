package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get method
		method := c.Request.Method

		// Log the request
		if query != "" {
			path = path + "?" + query
		}

		// Use structured logging
		logger.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)

		// Also log to console for development
		if gin.Mode() == gin.DebugMode {
			fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	}
}

// GetLogger returns the shared logger instance
func GetLogger() *zap.Logger {
	return logger
}
