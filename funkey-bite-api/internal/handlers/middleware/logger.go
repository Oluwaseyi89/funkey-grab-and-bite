package middleware

import (
	"fmt"
	"net/url"
	"strings"
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
		start := time.Now()
		path := c.Request.URL.Path
		query := redactQueryString(c.Request.URL.RawQuery)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()

		clientIP := c.ClientIP()

		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)

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

func redactQueryString(rawQuery string) string {
	if strings.TrimSpace(rawQuery) == "" {
		return ""
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return ""
	}

	sensitiveKeys := map[string]struct{}{
		"token":         {},
		"access_token":  {},
		"refresh_token": {},
		"authorization": {},
	}

	for key := range values {
		if _, isSensitive := sensitiveKeys[strings.ToLower(key)]; isSensitive {
			values.Set(key, "[REDACTED]")
		}
	}

	return values.Encode()
}

func GetLogger() *zap.Logger {
	return logger
}
