package v1

import (
	"errors"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers/middleware"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type AuthHandler struct {
	authService    services.AuthService
	userService    services.UserService
	accountService services.AccountService
	validate       *validator.Validate
}

func NewAuthHandler(authService services.AuthService, userService services.UserService, accountService services.AccountService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		userService:    userService,
		accountService: accountService,
		validate:       validator.New(),
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

	if !utils.ValidatePasswordStrength(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters with uppercase, lowercase, and numbers",
		})
		return
	}

	result, err := h.accountService.Register(middleware.RequestContext(c), req)
	if err != nil {
		respondAccountError(c, err)
		return
	}
	if result.VerificationSentTo != "" {
		// The number has guest orders; they are linked once the emailed
		// code is confirmed.
		c.JSON(http.StatusAccepted, gin.H{"verificationRequired": true, "sentTo": result.VerificationSentTo})
		return
	}

	c.JSON(http.StatusCreated, result.Auth)
}

// respondAccountError explains a registration or password reset problem.
func respondAccountError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrAccountExists), errors.Is(err, services.ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrInvalidCode), errors.Is(err, services.ErrCannotVerifyPhone):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Printf("Account request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again."})
	}
}

// ConfirmRegistration completes a registration for a number with guest
// orders, linking them to the new account
// @Summary Confirm registration with the emailed code
// @Tags auth
// @Accept json
// @Produce json
// @Param confirmation body models.RegistrationConfirmation true "Registration details and code"
// @Router /auth/register/confirm [post]
func (h *AuthHandler) ConfirmRegistration(c *gin.Context) {
	var req models.RegistrationConfirmation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter your details and the 6-digit code from the email"})
		return
	}
	if !utils.ValidatePasswordStrength(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters with uppercase, lowercase, and numbers"})
		return
	}

	auth, err := h.accountService.ConfirmRegistration(middleware.RequestContext(c), req.UserRegistration, req.Code)
	if err != nil {
		respondAccountError(c, err)
		return
	}
	c.JSON(http.StatusCreated, auth)
}

// ForgotPassword emails a reset code to the account with this email. The
// answer is the same whether or not one exists.
// @Summary Request a password reset code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.PasswordResetRequest true "Account email"
// @Router /auth/password/forgot [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid email"})
		return
	}
	if err := h.accountService.RequestPasswordReset(middleware.RequestContext(c), req.Email); err != nil {
		log.Printf("Password reset request failed: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "If an account uses that email, we've sent it a code."})
}

// ResetPassword sets a new password with an emailed code
// @Summary Reset a password
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body models.PasswordReset true "Email, code and new password"
// @Router /auth/password/reset [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.PasswordReset
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter your email, the 6-digit code and a new password"})
		return
	}
	if !utils.ValidatePasswordStrength(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters with uppercase, lowercase, and numbers"})
		return
	}
	if err := h.accountService.ResetPassword(middleware.RequestContext(c), req.Email, req.Code, req.Password); err != nil {
		respondAccountError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Your password has been changed. You can log in now."})
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

	// Intentionally return a generic message to prevent phone/email enumeration.
	c.JSON(http.StatusOK, gin.H{
		"message": "If an account exists for the provided details, continue with login.",
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

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

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body models.ProfileUpdate true "Profile updates"
// @Success 200 {object} models.User
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var updates models.ProfileUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the updates
	if err := h.validate.Struct(updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateProfile(userID.(int), &updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": updatedUser,
	})
}

// ChangePassword changes user password
// @Summary Change password
// @Description Change authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]string
// @Router /auth/password [patch]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ChangePassword(userID.(int), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
}
