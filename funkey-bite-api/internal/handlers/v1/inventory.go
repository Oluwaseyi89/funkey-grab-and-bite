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

type InventoryHandler struct {
	inventoryService services.InventoryService
	validate         *validator.Validate
}

func NewInventoryHandler(inventoryService services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		validate:         validator.New(),
	}
}

// GetInventory gets all inventory items
// @Summary Get all inventory items
// @Description Get all inventory items
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.InventoryItem
// @Router /admin/inventory [get]
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	items, err := h.inventoryService.GetAllInventory()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"INVENTORY_FETCH_FAILED", "Failed to fetch inventory", err.Error())
		return
	}

	handlers.Success(c, items)
}

// GetInventoryItem gets an inventory item by ID
// @Summary Get inventory item by ID
// @Description Get a specific inventory item by ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory Item ID"
// @Success 200 {object} models.InventoryItem
// @Router /admin/inventory/{id} [get]
func (h *InventoryHandler) GetInventoryItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid inventory item ID")
		return
	}

	item, err := h.inventoryService.GetInventoryItem(id)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"INVENTORY_FETCH_FAILED", "Failed to fetch inventory item", err.Error())
		return
	}

	if item == nil {
		handlers.Error(c, http.StatusNotFound, "NOT_FOUND", "Inventory item not found")
		return
	}

	handlers.Success(c, item)
}

// GetInventoryByMenuItem gets inventory by menu item ID
// @Summary Get inventory by menu item ID
// @Description Get inventory item by menu item ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param menuItemId path int true "Menu Item ID"
// @Success 200 {object} models.InventoryItem
// @Router /admin/inventory/menu-item/{menuItemId} [get]
func (h *InventoryHandler) GetInventoryByMenuItem(c *gin.Context) {
	menuItemID, err := strconv.Atoi(c.Param("menuItemId"))
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid menu item ID")
		return
	}

	item, err := h.inventoryService.GetInventoryByMenuItemID(menuItemID)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"INVENTORY_FETCH_FAILED", "Failed to fetch inventory item", err.Error())
		return
	}

	if item == nil {
		handlers.Error(c, http.StatusNotFound, "NOT_FOUND", "Inventory item not found for this menu item")
		return
	}

	handlers.Success(c, item)
}

// GetLowStock gets low stock inventory items
// @Summary Get low stock items
// @Description Get inventory items below reorder point
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.InventoryItem
// @Router /admin/inventory/low-stock [get]
func (h *InventoryHandler) GetLowStock(c *gin.Context) {
	items, err := h.inventoryService.GetLowStock()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"INVENTORY_FETCH_FAILED", "Failed to fetch low stock items", err.Error())
		return
	}

	handlers.Success(c, items)
}

// UpdateStock updates inventory stock
// @Summary Update inventory stock
// @Description Update inventory stock level
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param update body models.InventoryUpdate true "Stock update details"
// @Success 200 {object} models.InventoryItem
// @Router /admin/inventory/stock [patch]
func (h *InventoryHandler) UpdateStock(c *gin.Context) {
	var update models.InventoryUpdate

	if err := c.ShouldBindJSON(&update); err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.validate.Struct(update); err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"VALIDATION_FAILED", "Validation failed", err.Error())
		return
	}

	updatedItem, err := h.inventoryService.UpdateStock(&update)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"STOCK_UPDATE_FAILED", "Failed to update stock", err.Error())
		return
	}

	handlers.Success(c, updatedItem)
}

// CreateInventoryItem creates a new inventory item
// @Summary Create inventory item
// @Description Create a new inventory item
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body models.InventoryItem true "Inventory item details"
// @Success 201 {object} models.InventoryItem
// @Router /admin/inventory [post]
func (h *InventoryHandler) CreateInventoryItem(c *gin.Context) {
	var item models.InventoryItem

	if err := c.ShouldBindJSON(&item); err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Basic validation
	if item.MenuItemID <= 0 {
		handlers.Error(c, http.StatusBadRequest, "INVALID_MENU_ITEM", "Menu item ID is required")
		return
	}
	if item.Name == "" {
		handlers.Error(c, http.StatusBadRequest, "INVALID_NAME", "Item name is required")
		return
	}
	if item.CurrentStock < 0 {
		handlers.Error(c, http.StatusBadRequest, "INVALID_STOCK", "Stock cannot be negative")
		return
	}

	createdItem, err := h.inventoryService.CreateInventoryItem(&item)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"INVENTORY_CREATE_FAILED", "Failed to create inventory item", err.Error())
		return
	}

	handlers.Created(c, createdItem)
}

// GetAlerts gets inventory alerts
// @Summary Get inventory alerts
// @Description Get active or resolved inventory alerts
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resolved query bool false "Filter by resolved status"
// @Success 200 {array} models.InventoryAlert
// @Router /admin/inventory/alerts [get]
func (h *InventoryHandler) GetAlerts(c *gin.Context) {
	resolvedStr := c.Query("resolved")
	var resolved bool
	if resolvedStr != "" {
		var err error
		resolved, err = strconv.ParseBool(resolvedStr)
		if err != nil {
			handlers.Error(c, http.StatusBadRequest, "INVALID_PARAM", "Invalid resolved parameter")
			return
		}
	}

	alerts, err := h.inventoryService.GetAlerts(resolved)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"ALERTS_FETCH_FAILED", "Failed to fetch alerts", err.Error())
		return
	}

	handlers.Success(c, alerts)
}

// ResolveAlert resolves an inventory alert
// @Summary Resolve inventory alert
// @Description Mark an inventory alert as resolved
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Alert ID"
// @Success 200 {object} map[string]string
// @Router /admin/inventory/alerts/{id}/resolve [patch]
func (h *InventoryHandler) ResolveAlert(c *gin.Context) {
	alertID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid alert ID")
		return
	}

	err = h.inventoryService.ResolveAlert(alertID)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"ALERT_RESOLVE_FAILED", "Failed to resolve alert", err.Error())
		return
	}

	handlers.Success(c, gin.H{
		"message": "Alert resolved successfully",
		"alertId": alertID,
	})
}

// CheckAvailability checks if a menu item is available
// @Summary Check item availability
// @Description Check if a menu item has sufficient stock
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param menuItemId query int true "Menu Item ID"
// @Param quantity query int true "Requested quantity"
// @Success 200 {object} map[string]interface{}
// @Router /admin/inventory/check [get]
func (h *InventoryHandler) CheckAvailability(c *gin.Context) {
	menuItemID, err := strconv.Atoi(c.Query("menuItemId"))
	if err != nil || menuItemID <= 0 {
		handlers.Error(c, http.StatusBadRequest, "INVALID_MENU_ITEM", "Valid menu item ID is required")
		return
	}

	quantity, err := strconv.Atoi(c.Query("quantity"))
	if err != nil || quantity <= 0 {
		handlers.Error(c, http.StatusBadRequest, "INVALID_QUANTITY", "Valid quantity is required")
		return
	}

	isAvailable, message := h.inventoryService.CheckAvailability(menuItemID, quantity)

	handlers.Success(c, gin.H{
		"menuItemId":  menuItemID,
		"quantity":    quantity,
		"isAvailable": isAvailable,
		"message":     message,
	})
}

// RestockItem restocks an inventory item
// @Summary Restock inventory item
// @Description Add stock to an inventory item
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param restock body RestockRequest true "Restock details"
// @Success 200 {object} models.InventoryItem
// @Router /admin/inventory/restock [post]
func (h *InventoryHandler) RestockItem(c *gin.Context) {
	var req RestockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.MenuItemID <= 0 {
		handlers.Error(c, http.StatusBadRequest, "INVALID_MENU_ITEM", "Menu item ID is required")
		return
	}
	if req.Quantity <= 0 {
		handlers.Error(c, http.StatusBadRequest, "INVALID_QUANTITY", "Quantity must be positive")
		return
	}
	if req.Reason == "" {
		handlers.Error(c, http.StatusBadRequest, "MISSING_REASON", "Reason is required")
		return
	}

	err := h.inventoryService.RestockItem(req.MenuItemID, req.Quantity, req.Reason)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusBadRequest,
			"RESTOCK_FAILED", "Failed to restock item", err.Error())
		return
	}

	// Get updated item
	updatedItem, err := h.inventoryService.GetInventoryByMenuItemID(req.MenuItemID)
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"INVENTORY_FETCH_FAILED", "Failed to fetch updated item", err.Error())
		return
	}

	handlers.Success(c, updatedItem)
}

// GetDashboard gets inventory dashboard statistics
// @Summary Get inventory dashboard
// @Description Get inventory dashboard statistics
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /admin/inventory/dashboard [get]
func (h *InventoryHandler) GetDashboard(c *gin.Context) {
	dashboard, err := h.inventoryService.GetInventoryDashboard()
	if err != nil {
		handlers.ErrorWithDetails(c, http.StatusInternalServerError,
			"DASHBOARD_FETCH_FAILED", "Failed to fetch dashboard", err.Error())
		return
	}

	handlers.Success(c, dashboard)
}

// RestockRequest struct for restock endpoint
type RestockRequest struct {
	MenuItemID int    `json:"menuItemId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Reason     string `json:"reason" binding:"required"`
}
