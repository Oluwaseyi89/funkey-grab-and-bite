package v1

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/realtime"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type CateringHandler struct {
	cateringService services.CateringService
	validate        *validator.Validate
}

func NewCateringHandler(cateringService services.CateringService) *CateringHandler {
	validate := validator.New()
	// Report fields by their JSON names so messages match what clients send.
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CateringHandler{
		cateringService: cateringService,
		validate:        validate,
	}
}

// cateringFieldLabels names request and package fields in plain words.
var cateringFieldLabels = map[string]string{
	"requestType": "Request type", "packageId": "Package", "contactName": "Your name",
	"contactPhone": "Phone number", "contactEmail": "Email address", "organization": "Organisation",
	"eventType": "Event type", "eventTypeOther": "Event description", "eventName": "Event name",
	"eventDate": "Event date", "eventTime": "Start time", "durationHours": "Duration",
	"guestCount": "Number of guests", "venueName": "Venue name", "venueAddress": "Venue address",
	"venueArea": "Area or city", "venueState": "State", "venueNotes": "Venue notes",
	"serviceStyle": "Service style", "selectedItemIds": "Chosen dishes", "menuNotes": "Menu notes",
	"dietaryRequirements": "Dietary requirements", "dietaryNotes": "Dietary notes", "budget": "Budget",
	"specialRequests": "Special requests", "name": "Package name", "description": "Description",
	"pricePerPerson": "Price per person", "minGuests": "Minimum guests", "maxGuests": "Maximum guests",
	"minNoticeDays": "Minimum notice", "inclusions": "What's included", "dishChoices": "Dish choices",
	"menuItemIds": "Menu items", "displayOrder": "Display order", "quotedAmount": "Quoted amount",
	"adminNotes": "Notes", "reference": "Reference", "phone": "Phone number",
}

// validationMessage turns the first validation failure into a sentence,
// naming the field with labels.
func validationMessage(err error, labels map[string]string) string {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) || len(errs) == 0 {
		return "Please check the details and try again"
	}

	fieldErr := errs[0]
	field := fieldErr.Field()
	if i := strings.Index(field, "["); i > 0 {
		field = field[:i]
	}
	label, ok := labels[field]
	if !ok {
		label = field
	}

	switch fieldErr.Tag() {
	case "required":
		return label + " is required"
	case "email":
		return "Enter a valid email address"
	case "oneof":
		return label + " has an option we don't offer"
	case "min", "gt":
		return fmt.Sprintf("%s is too small or too short", label)
	case "max", "lte":
		return fmt.Sprintf("%s is too large or too long", label)
	default:
		return label + " is invalid"
	}
}

// respondCateringError maps service errors to responses.
func respondCateringError(c *gin.Context, err error, fallback string) {
	var invalid *services.CateringValidationError
	switch {
	case errors.As(err, &invalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": invalid.Message})
	case errors.Is(err, services.ErrCateringRequestNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Catering request not found"})
	case errors.Is(err, services.ErrCateringPackageNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Catering package not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallback})
	}
}

// customerView hides what only the business should see.
func customerView(request *models.CateringRequest) *models.CateringRequest {
	view := *request
	view.AdminNotes = nil
	view.UserID = nil
	return &view
}

// CreateRequest creates a catering request for a package or a custom event
// @Summary Create catering request
// @Description Request catering for a package or a custom event; returns the request with its reference
// @Tags catering
// @Accept json
// @Produce json
// @Param request body models.CateringRequestInput true "Catering request details"
// @Success 201 {object} models.CateringRequest
// @Router /catering/requests [post]
func (h *CateringHandler) CreateRequest(c *gin.Context) {
	var input models.CateringRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err, cateringFieldLabels)})
		return
	}

	var userID *int
	if uid, exists := c.Get("user_id"); exists {
		id := uid.(int)
		userID = &id
	}

	request, err := h.cateringService.CreateRequest(input, userID)
	if err != nil {
		respondCateringError(c, err, "Failed to submit catering request")
		return
	}

	realtime.GlobalBroadcaster.Broadcast("new_catering_request", request)
	c.JSON(http.StatusCreated, customerView(request))
}

// LookupRequest lets a customer check their request by reference and phone
// @Summary Look up catering request
// @Description Get a catering request by its reference and the phone number it was made with
// @Tags catering
// @Accept json
// @Produce json
// @Param lookup body models.CateringLookupRequest true "Reference and phone"
// @Success 200 {object} models.CateringRequest
// @Router /catering/requests/lookup [post]
func (h *CateringHandler) LookupRequest(c *gin.Context) {
	var input models.CateringLookupRequest
	if err := c.ShouldBindJSON(&input); err != nil || h.validate.Struct(input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reference and phone number are required"})
		return
	}

	request, err := h.cateringService.LookupRequest(input.Reference, input.Phone)
	if err != nil {
		respondCateringError(c, err, "Failed to fetch catering request")
		return
	}
	c.JSON(http.StatusOK, customerView(request))
}

// GetRequest gets a catering request by ID
// @Summary Get catering request
// @Description Get a specific catering request by ID (its owner or an admin)
// @Tags catering
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Success 200 {object} models.CateringRequest
// @Router /admin/catering/requests/{id} [get]
func (h *CateringHandler) GetRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	request, err := h.cateringService.GetRequestByID(id)
	if err != nil {
		respondCateringError(c, err, "Failed to fetch catering request")
		return
	}

	if isAdmin, _ := c.Get("is_admin"); isAdmin == true {
		c.JSON(http.StatusOK, request)
		return
	}
	userID, _ := c.Get("user_id")
	if request.UserID == nil || *request.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	c.JSON(http.StatusOK, customerView(request))
}

// GetUserRequests gets the signed-in customer's catering requests
// @Summary Get my catering requests
// @Tags catering
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/catering-requests [get]
func (h *CateringHandler) GetUserRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	requests, err := h.cateringService.GetUserRequests(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch catering requests"})
		return
	}

	views := make([]*models.CateringRequest, 0, len(requests))
	for i := range requests {
		views = append(views, customerView(&requests[i]))
	}
	c.JSON(http.StatusOK, gin.H{"requests": views, "count": len(views)})
}

// GetAllRequests gets all catering requests (admin only)
// @Summary Get all catering requests
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /admin/catering/requests [get]
func (h *CateringHandler) GetAllRequests(c *gin.Context) {
	requests, err := h.cateringService.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch catering requests"})
		return
	}
	if requests == nil {
		requests = []models.CateringRequest{}
	}
	c.JSON(http.StatusOK, gin.H{"requests": requests, "count": len(requests)})
}

// UpdateRequest records the admin's follow-up on a request
// @Summary Update catering request
// @Description Change a request's status, quoted amount or internal notes (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Param update body models.CateringRequestUpdate true "Status, quote and notes"
// @Success 200 {object} models.CateringRequest
// @Router /admin/catering/requests/{id} [patch]
func (h *CateringHandler) UpdateRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var update models.CateringRequestUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err, cateringFieldLabels)})
		return
	}

	request, err := h.cateringService.UpdateRequest(id, &update)
	if err != nil {
		respondCateringError(c, err, "Failed to update catering request")
		return
	}
	c.JSON(http.StatusOK, request)
}

// UpdateRequestStatus updates a catering request's status (admin only)
// @Summary Update catering request status
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Success 200 {object} models.CateringRequest
// @Router /admin/catering/requests/{id}/status [patch]
func (h *CateringHandler) UpdateRequestStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	request, err := h.cateringService.UpdateRequest(id, &models.CateringRequestUpdate{Status: &body.Status})
	if err != nil {
		respondCateringError(c, err, "Failed to update status")
		return
	}
	c.JSON(http.StatusOK, request)
}

// ListPackages lists the catering packages customers can choose
// @Summary List catering packages
// @Tags catering
// @Produce json
// @Success 200 {array} models.CateringPackage
// @Router /catering/packages [get]
func (h *CateringHandler) ListPackages(c *gin.Context) {
	h.listPackages(c, true)
}

// ListAllPackages lists every catering package, including switched-off ones
// @Summary List all catering packages (admin)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CateringPackage
// @Router /admin/catering/packages [get]
func (h *CateringHandler) ListAllPackages(c *gin.Context) {
	h.listPackages(c, false)
}

func (h *CateringHandler) listPackages(c *gin.Context, activeOnly bool) {
	packages, err := h.cateringService.ListPackages(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch catering packages"})
		return
	}
	c.JSON(http.StatusOK, packages)
}

// GetPackage gets one catering package (admin)
// @Summary Get catering package
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Success 200 {object} models.CateringPackage
// @Router /admin/catering/packages/{id} [get]
func (h *CateringHandler) GetPackage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID"})
		return
	}
	pkg, err := h.cateringService.GetPackage(id)
	if err != nil {
		respondCateringError(c, err, "Failed to fetch catering package")
		return
	}
	c.JSON(http.StatusOK, pkg)
}

// CreatePackage creates a catering package (admin)
// @Summary Create catering package
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param package body models.CateringPackageInput true "Package details"
// @Success 201 {object} models.CateringPackage
// @Router /admin/catering/packages [post]
func (h *CateringHandler) CreatePackage(c *gin.Context) {
	input, ok := h.bindPackage(c)
	if !ok {
		return
	}
	pkg, err := h.cateringService.CreatePackage(input)
	if err != nil {
		respondCateringError(c, err, "Failed to create catering package")
		return
	}
	c.JSON(http.StatusCreated, pkg)
}

// UpdatePackage replaces a catering package (admin)
// @Summary Update catering package
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Package ID"
// @Param package body models.CateringPackageInput true "Package details"
// @Success 200 {object} models.CateringPackage
// @Router /admin/catering/packages/{id} [put]
func (h *CateringHandler) UpdatePackage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID"})
		return
	}
	input, ok := h.bindPackage(c)
	if !ok {
		return
	}
	pkg, err := h.cateringService.UpdatePackage(id, input)
	if err != nil {
		respondCateringError(c, err, "Failed to update catering package")
		return
	}
	c.JSON(http.StatusOK, pkg)
}

func (h *CateringHandler) bindPackage(c *gin.Context) (*models.CateringPackageInput, bool) {
	var input models.CateringPackageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return nil, false
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err, cateringFieldLabels)})
		return nil, false
	}
	return &input, true
}
