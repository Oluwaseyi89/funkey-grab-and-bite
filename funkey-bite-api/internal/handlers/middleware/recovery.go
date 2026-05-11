package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logger := GetLogger()
				requestID := GetRequestID(c)
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(stack)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("request_id", requestID),
				)

				if gin.Mode() == gin.DebugMode {
					fmt.Printf("Panic: %v\n", err)
					fmt.Printf("Stack: %s\n", stack)
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal server error",
					"code":      "INTERNAL_SERVER_ERROR",
					"requestId": requestID,
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
