package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

// RequestIDMiddleware ensures each request has a correlation id for logs and responses.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if value, exists := c.Get(requestIDKey); exists {
		if id, ok := value.(string); ok && strings.TrimSpace(id) != "" {
			return id
		}
	}

	requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
	if requestID != "" {
		return requestID
	}

	return "unknown"
}

func generateRequestID() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(buffer)
}
