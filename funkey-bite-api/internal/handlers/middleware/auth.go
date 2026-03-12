package middleware

import (
	"net/http"
	"strings"

	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Bearer token required",
				"code":  "INVALID_TOKEN_FORMAT",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			if claims, err := utils.ValidateToken(tokenString); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("phone", claims.Phone)
			}
		}
		c.Next()
	}
}
