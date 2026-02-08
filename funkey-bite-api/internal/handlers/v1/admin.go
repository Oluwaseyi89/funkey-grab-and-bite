package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type AdminHandler struct {
	adminService services.AdminService
	validate     *validator.Validate
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		validate:     validator.New(),
	}
}

// GetDashboardStats returns admin dashboard statistics
// @Summary Get dashboard statistics
// @Description Get admin dashboard statistics
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.AdminStats
// @Router /admin/dashboard/stats [get]
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get dashboard statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSalesReport returns sales report
// @Summary Get sales report
// @Description Get sales report for a date range
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "From date (YYYY-MM-DD)"
// @Param to query string true "To date (YYYY-MM-DD)"
// @Success 200 {array} models.SalesReport
// @Router /admin/reports/sales [get]
func (h *AdminHandler) GetSalesReport(c *gin.Context) {
	fromDate := c.Query("from")
	toDate := c.Query("to")

	if fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to dates are required",
		})
		return
	}

	report, err := h.adminService.GetSalesReport(fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetAllOrders returns all orders with pagination
// @Summary Get all orders
// @Description Get all orders with pagination and filtering
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Router /admin/orders [get]
func (h *AdminHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	orders, total, err := h.adminService.GetAllOrders(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// UpdateOrderStatus updates order status
// @Summary Update order status
// @Description Update the status of an order
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param status body string true "New status"
// @Success 200 {object} map[string]string
// @Router /admin/orders/{id}/status [patch]
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending confirmed preparing ready completed cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = h.adminService.UpdateOrderStatus(orderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"orderId": orderID,
		"status":  req.Status,
	})
}

// GetAllUsers returns all users with pagination
// @Summary Get all users
// @Description Get all users with pagination
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /admin/users [get]
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.adminService.GetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// UpdateUserStatus updates user active status
// @Summary Update user status
// @Description Update user active/inactive status
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param status body bool true "Active status"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id}/status [patch]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req struct {
		IsActive bool `json:"isActive" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = h.adminService.UpdateUserStatus(userID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User status updated successfully",
		"userId":   userID,
		"isActive": req.IsActive,
	})
}

// CreateMenuItem creates a new menu item
// @Summary Create menu item
// @Description Create a new menu item
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body models.MenuItem true "Menu item details"
// @Success 201 {object} models.MenuItem
// @Router /admin/menu/items [post]
func (h *AdminHandler) CreateMenuItem(c *gin.Context) {
	var item models.MenuItem

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	createdItem, err := h.adminService.CreateMenuItem(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdItem)
}

// UpdateMenuItem updates a menu item
// @Summary Update menu item
// @Description Update an existing menu item
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu Item ID"
// @Param item body models.MenuItem true "Updated menu item details"
// @Success 200 {object} models.MenuItem
// @Router /admin/menu/items/{id} [put]
func (h *AdminHandler) UpdateMenuItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid menu item ID",
		})
		return
	}

	var item models.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	item.ID = itemID

	err = h.adminService.UpdateMenuItem(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return updated item
	updatedItem, err := h.adminService.GetMenuItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated item",
		})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

// DeleteMenuItem deletes a menu item
// @Summary Delete menu item
// @Description Delete a menu item
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu Item ID"
// @Success 200 {object} map[string]string
// @Router /admin/menu/items/{id} [delete]
func (h *AdminHandler) DeleteMenuItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid menu item ID",
		})
		return
	}

	err = h.adminService.DeleteMenuItem(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu item deleted successfully",
		"itemId":  itemID,
	})
}
