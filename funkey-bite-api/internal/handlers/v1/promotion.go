package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type PromotionHandler struct {
	promotionService services.PromotionService
	validate         *validator.Validate
}

func NewPromotionHandler(promotionService services.PromotionService) *PromotionHandler {
	return &PromotionHandler{
		promotionService: promotionService,
		validate:         validator.New(),
	}
}

// CreatePromotion creates a new promotion
// @Summary Create promotion
// @Description Create a new promotion
// @Tags promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param promotion body models.PromotionCreate true "Promotion details"
// @Success 201 {object} models.Promotion
// @Router /admin/promotions [post]
func (h *PromotionHandler) CreatePromotion(c *gin.Context) {
	var create models.PromotionCreate

	if err := c.ShouldBindJSON(&create); err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.validate.Struct(create); err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", err.Error())
		return
	}

	promotion, err := h.promotionService.CreatePromotion(&create)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest, "PROMOTION_CREATE_FAILED", "Failed to create promotion", err.Error())
		return
	}

	handlers.Created(c, promotion)
}

// GetPromotions gets all promotions with pagination
// @Summary Get all promotions
// @Description Get all promotions with pagination
// @Tags promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status (active/inactive)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/promotions [get]
func (h *PromotionHandler) GetPromotions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	promotions, total, err := h.promotionService.GetAllPromotions(page, limit, status)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"PROMOTIONS_FETCH_FAILED", "Failed to fetch promotions", err.Error())
		return
	}

	handlers.Paginated(c, promotions, page, limit, total)
}

// GetPromotion gets a promotion by ID
// @Summary Get promotion by ID
// @Description Get a specific promotion by ID
// @Tags promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Promotion ID"
// @Success 200 {object} models.Promotion
// @Router /admin/promotions/{id} [get]
func (h *PromotionHandler) GetPromotion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid promotion ID")
		return
	}

	promotion, err := h.promotionService.GetPromotionByID(id)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"PROMOTION_FETCH_FAILED", "Failed to fetch promotion", err.Error())
		return
	}

	if promotion == nil {
		handlers.Error(c, http.StatusNotFound, "NOT_FOUND", "Promotion not found")
		return
	}

	handlers.Success(c, promotion)
}

// UpdatePromotion updates a promotion
// @Summary Update promotion
// @Description Update an existing promotion
// @Tags promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Promotion ID"
// @Param promotion body models.PromotionUpdate true "Promotion updates"
// @Success 200 {object} models.Promotion
// @Router /admin/promotions/{id} [put]
func (h *PromotionHandler) UpdatePromotion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid promotion ID")
		return
	}

	var updates models.PromotionUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.validate.Struct(updates); err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", err.Error())
		return
	}

	promotion, err := h.promotionService.UpdatePromotion(id, &updates)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"PROMOTION_UPDATE_FAILED", "Failed to update promotion", err.Error())
		return
	}

	handlers.Success(c, promotion)
}

// DeletePromotion deletes a promotion
// @Summary Delete promotion
// @Description Delete a promotion
// @Tags promotions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Promotion ID"
// @Success 200 {object} map[string]string
// @Router /admin/promotions/{id} [delete]
func (h *PromotionHandler) DeletePromotion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid promotion ID")
		return
	}

	err = h.promotionService.DeletePromotion(id)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"PROMOTION_DELETE_FAILED", "Failed to delete promotion", err.Error())
		return
	}

	handlers.Success(c, gin.H{"message": "Promotion deleted successfully", "id": id})
}

// ValidatePromotion validates a promotion code
// @Summary Validate promotion
// @Description Validate a promotion code
// @Tags promotions
// @Accept json
// @Produce json
// @Param code query string true "Promotion code"
// @Param amount query number false "Order amount for validation"
// @Success 200 {object} models.PromotionValidation
// @Router /promotions/validate [get]
func (h *PromotionHandler) ValidatePromotion(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		handlers.Error(c, http.StatusBadRequest, "MISSING_CODE", "Promotion code is required")
		return
	}

	amount, _ := strconv.ParseFloat(c.Query("amount"), 64)
	if amount < 0 {
		amount = 0
	}

	validation, err := h.promotionService.ValidatePromotion(code, amount, nil)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"VALIDATION_FAILED", "Failed to validate promotion", err.Error())
		return
	}

	handlers.Success(c, validation)
}

// GetActivePromotions gets active promotions for public display
// @Summary Get active promotions
// @Description Get active promotions for public display
// @Tags promotions
// @Accept json
// @Produce json
// @Success 200 {array} models.Promotion
// @Router /promotions/active [get]
func (h *PromotionHandler) GetActivePromotions(c *gin.Context) {
	promotions, err := h.promotionService.GetActivePromotions()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"PROMOTIONS_FETCH_FAILED", "Failed to fetch promotions", err.Error())
		return
	}

	handlers.Success(c, promotions)
}
