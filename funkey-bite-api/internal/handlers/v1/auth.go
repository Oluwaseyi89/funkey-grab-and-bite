// internal/handlers/v1/auth.go
package v1

import (
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
	userService services.UserService // Add this
	validate    *validator.Validate
}

func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService, // Add this
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegistration

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password strength
	if !utils.ValidatePasswordStrength(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters with uppercase, lowercase, and numbers",
		})
		return
	}

	authResponse, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, authResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := h.authService.Login(req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

func (h *AuthHandler) CheckUser(c *gin.Context) {
	phone := c.Query("phone")
	email := c.Query("email")

	if phone == "" && email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone or email required"})
		return
	}

	user, exists, err := h.authService.CheckUserExists(phone, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": exists,
		"user":   user,
	})
}

// func (h *AuthHandler) GetProfile(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
// 		return
// 	}

// 	// Fetch user from service
// 	// user, err := h.userService.GetByID(userID.(int))
// 	// Handle response...
// }

// Then use the proper GetProfile method:
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Cast userID to int
	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Use userService to get profile
	user, err := h.userService.GetByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile: " + err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
