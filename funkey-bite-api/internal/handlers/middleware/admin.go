package middleware

import (
	"database/sql"
	"net/http"

	// "funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the authenticated user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get database from context or use a global connection
		dbInterface, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
			c.Abort()
			return
		}

		db, ok := dbInterface.(*sql.DB)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
			c.Abort()
			return
		}

		// Check if user is admin
		isAdmin, err := isUserAdmin(db, userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin status"})
			c.Abort()
			return
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isUserAdmin checks if a user exists in the admin_users table
func isUserAdmin(db *sql.DB, userID int) (bool, error) {
	// You can implement this in multiple ways:
	// Option 1: Check admin_users table
	query := `SELECT EXISTS(SELECT 1 FROM admin_users WHERE id = $1 AND is_active = true)`
	var exists bool
	err := db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

	// Option 2: Check user role in users table (if you have role field)
	// query := `SELECT role FROM users WHERE id = $1 AND is_active = true`
	// var role string
	// err := db.QueryRow(query, userID).Scan(&role)
	// if err != nil {
	//     return false, err
	// }
	// return role == "admin", nil
}

// AdminAuthMiddleware combines authentication and admin check
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First authenticate
		AuthMiddleware()(c)

		// If authentication passed, check admin status
		if c.IsAborted() {
			return
		}

		// Then check admin
		AdminMiddleware()(c)
	}
}
