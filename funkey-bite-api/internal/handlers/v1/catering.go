package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type CateringHandler struct {
	cateringService services.CateringService
	validate        *validator.Validate
}

func NewCateringHandler(cateringService services.CateringService) *CateringHandler {
	return &CateringHandler{
		cateringService: cateringService,
		validate:        validator.New(),
	}
}

// CreateRequest creates a new catering request
// @Summary Create catering request
// @Description Create a new catering request
// @Tags catering
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CateringRequestInput true "Catering request details"
// @Success 201 {object} models.CateringRequest
// @Router /catering/requests [post]
func (h *CateringHandler) CreateRequest(c *gin.Context) {
	var input models.CateringRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Get user ID from context if authenticated
	var userID *int
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(int)
		userID = &id
	}

	request, err := h.cateringService.CreateRequest(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create catering request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// GetRequest gets a catering request by ID
// @Summary Get catering request
// @Description Get a specific catering request by ID
// @Tags catering
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Success 200 {object} models.CateringRequest
// @Router /catering/requests/{id} [get]
func (h *CateringHandler) GetRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request ID",
		})
		return
	}

	request, err := h.cateringService.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Catering request not found",
		})
		return
	}

	userID, _ := c.Get("user_id")
	if request.UserID != nil && *request.UserID != userID {
		isAdmin, _ := c.Get("is_admin")
		if isAdmin == nil || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
	}

	c.JSON(http.StatusOK, request)
}

// GetUserRequests gets all catering requests for the authenticated user
// @Summary Get user catering requests
// @Description Get all catering requests for the authenticated user
// @Tags catering
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CateringRequest
// @Router /catering/requests/user [get]
func (h *CateringHandler) GetUserRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	requests, err := h.cateringService.GetUserRequests(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch catering requests",
		})
		return
	}
	if requests == nil {
		requests = []models.CateringRequest{}
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"count":    len(requests),
	})
}

// GetPackages gets available catering packages
// @Summary Get catering packages
// @Description Get all available catering packages
// @Tags catering
// @Accept json
// @Produce json
// @Success 200 {array} models.CateringPackage
// @Router /catering/packages [get]
func (h *CateringHandler) GetPackages(c *gin.Context) {
	packages := h.cateringService.GetPackages()
	c.JSON(http.StatusOK, packages)
}

// UpdateRequestStatus updates a catering request status (admin only)
// @Summary Update catering request status
// @Description Update the status of a catering request (admin only)
// @Tags catering
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Param status body string true "New status"
// @Success 200 {object} map[string]string
// @Router /catering/requests/{id}/status [patch]
func (h *CateringHandler) UpdateRequestStatus(c *gin.Context) {
	isAdmin, _ := c.Get("is_admin")
	if isAdmin == nil || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request ID",
		})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = h.cateringService.UpdateRequestStatus(id, statusUpdate.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Status updated successfully",
		"requestId": id,
		"status":    statusUpdate.Status,
	})
}

// GetAllRequests gets all catering requests (admin only)
// @Summary Get all catering requests
// @Description Get all catering requests (admin only)
// @Tags catering
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CateringRequest
// @Router /catering/requests [get]
func (h *CateringHandler) GetAllRequests(c *gin.Context) {
	isAdmin, _ := c.Get("is_admin")
	if isAdmin == nil || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Admin access required",
		})
		return
	}

	requests, err := h.cateringService.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch catering requests",
		})
		return
	}
	if requests == nil {
		requests = []models.CateringRequest{}
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"count":    len(requests),
	})
}
