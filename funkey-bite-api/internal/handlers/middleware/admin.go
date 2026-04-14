package middleware

import (
	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminLookup interface {
	GetAdminUserByID(adminID int) (*models.AdminUser, error)
}

func AdminMiddleware(adminLookup AdminLookup) gin.HandlerFunc {
	return adminAccessMiddleware(adminLookup)
}

func AdminAuthMiddleware(adminLookup AdminLookup) gin.HandlerFunc {
	return adminAccessMiddleware(adminLookup)
}

func adminAccessMiddleware(adminLookup AdminLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if !claims.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		if adminLookup == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin authentication is not configured"})
			c.Abort()
			return
		}

		admin, err := adminLookup.GetAdminUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin access"})
			c.Abort()
			return
		}

		if admin == nil || !admin.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

var _ repository.IAdminRepository
