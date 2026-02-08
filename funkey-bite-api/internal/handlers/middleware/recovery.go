package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				stack := debug.Stack()
				logger := GetLogger()
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(stack)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Print to console in development
				if gin.Mode() == gin.DebugMode {
					fmt.Printf("Panic: %v\n", err)
					fmt.Printf("Stack: %s\n", stack)
				}

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"code":  "INTERNAL_SERVER_ERROR",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
