package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type SettingsHandler struct {
	settingsService services.SettingsService
	validate        *validator.Validate
}

func NewSettingsHandler(settingsService services.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		validate:        validator.New(),
	}
}

// GetSettings returns business settings
// @Summary Get business settings
// @Description Get business settings and configuration
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.BusinessSettings
// @Router /admin/settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingsService.GetSettings()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"SETTINGS_FETCH_FAILED", "Failed to fetch settings", err.Error())
		return
	}

	handlers.Success(c, settings)
}

// UpdateSettings updates business settings
// @Summary Update business settings
// @Description Update business settings and configuration
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body models.BusinessSettingsUpdate true "Settings updates"
// @Success 200 {object} models.BusinessSettings
// @Router /admin/settings [put]
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var updates models.BusinessSettingsUpdate

	if err := c.ShouldBindJSON(&updates); err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate the updates
	if err := h.validate.Struct(updates); err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"VALIDATION_FAILED", "Validation failed", err.Error())
		return
	}

	updatedSettings, err := h.settingsService.UpdateSettings(&updates)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"SETTINGS_UPDATE_FAILED", "Failed to update settings", err.Error())
		return
	}

	handlers.Success(c, updatedSettings)
}

// GetPublicSettings returns public business settings
// @Summary Get public business settings
// @Description Get public business settings (no auth required)
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /settings [get]
func (h *SettingsHandler) GetPublicSettings(c *gin.Context) {
	settings, err := h.settingsService.GetSettings()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"SETTINGS_FETCH_FAILED", "Failed to fetch settings", err.Error())
		return
	}

	// Return only public fields
	handlers.Success(c, gin.H{
		"businessName":   settings.BusinessName,
		"phoneNumber":    settings.PhoneNumber,
		"email":          settings.Email,
		"address":        settings.Address,
		"openingHours":   settings.OpeningHours,
		"deliveryFee":    settings.DeliveryFee,
		"minOrderAmount": settings.MinOrderAmount,
		"taxRate":        settings.TaxRate,
		"isDeliveryOpen": settings.IsDeliveryOpen,
		"isPickupOpen":   settings.IsPickupOpen,
	})
}

// GetOpeningHours returns business opening hours
// @Summary Get opening hours
// @Description Get business opening hours
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {array} models.BusinessHours
// @Router /settings/hours [get]
func (h *SettingsHandler) GetOpeningHours(c *gin.Context) {
	hours, err := h.settingsService.GetOpeningHours()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"HOURS_FETCH_FAILED", "Failed to fetch opening hours", err.Error())
		return
	}

	handlers.Success(c, hours)
}
