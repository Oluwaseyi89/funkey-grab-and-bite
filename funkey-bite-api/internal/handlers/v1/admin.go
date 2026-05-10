package v1

import (
	"database/sql"
	"errors"
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
	if report == nil {
		report = []models.SalesReport{}
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
	if orders == nil {
		orders = []models.Order{}
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
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
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
	if users == nil {
		users = []models.User{}
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
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
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

// GetMenuItems returns menu items for admin screens.
// @Summary Get menu items
// @Description Get menu items with optional pagination and filters
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param categoryId query int false "Filter by category ID"
// @Param query query string false "Search by name/description"
// @Success 200 {array} models.MenuItem
// @Router /admin/menu/items [get]
func (h *AdminHandler) GetMenuItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	query := c.Query("query")

	var categoryID *int
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		categoryID = &id
	}

	items, err := h.adminService.GetMenuItems(page, limit, categoryID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menu items"})
		return
	}
	if items == nil {
		items = []models.MenuItem{}
	}

	c.JSON(http.StatusOK, items)
}

// GetMenuItem returns a single menu item by ID.
// @Summary Get menu item
// @Description Get a menu item by ID for admin screens
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Router /admin/menu/items/{id} [get]
func (h *AdminHandler) GetMenuItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	item, err := h.adminService.GetMenuItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get menu item"})
		return
	}

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
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

// CreateMenuCategory creates a new menu category.
// @Summary Create menu category
// @Description Create a menu category for admin category management
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body models.MenuCategory true "Category details"
// @Success 201 {object} models.MenuCategory
// @Router /admin/menu/categories [post]
func (h *AdminHandler) CreateMenuCategory(c *gin.Context) {
	var category models.MenuCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	createdCategory, err := h.adminService.CreateMenuCategory(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

type updateCategoryRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	DisplayOrder *int    `json:"displayOrder"`
	IsActive     *bool   `json:"isActive"`
}

// UpdateMenuCategory updates an existing menu category.
// @Summary Update menu category
// @Description Update a menu category for admin category management
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param category body updateCategoryRequest true "Category updates"
// @Success 200 {object} models.MenuCategory
// @Router /admin/menu/categories/{id} [put]
func (h *AdminHandler) UpdateMenuCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	existingCategory, err := h.adminService.GetMenuCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		return
	}

	if existingCategory == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Name != nil {
		existingCategory.Name = *req.Name
	}
	if req.Description != nil {
		existingCategory.Description = *req.Description
	}
	if req.DisplayOrder != nil {
		existingCategory.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		existingCategory.IsActive = *req.IsActive
	}

	if err := h.adminService.UpdateMenuCategory(existingCategory); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingCategory)
}

// Add to the existing AdminHandler struct and methods...

// AdminLogin authenticates admin users
// @Summary Admin login
// @Description Authenticate admin user and return JWT token
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param credentials body AdminLoginRequest true "Admin credentials"
// @Success 200 {object} AdminLoginResponse
// @Router /admin/auth/login [post]
func (h *AdminHandler) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	admin, token, err := h.adminService.AdminLogin(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Clear sensitive data
	admin.PasswordHash = ""

	c.JSON(http.StatusOK, AdminLoginResponse{
		Admin: admin,
		Token: token,
	})
}

// AdminLogout logs out admin user
// @Summary Admin logout
// @Description Logout admin user (client-side token invalidation)
// @Tags admin-auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /admin/auth/logout [post]
func (h *AdminHandler) AdminLogout(c *gin.Context) {
	// For JWT tokens, logout is client-side
	// Could implement token blacklist here if needed
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetAdminUsers returns all admin users with pagination
// @Summary Get admin users
// @Description Get all admin users with pagination
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /admin/users/admins [get]
func (h *AdminHandler) GetAdminUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	admins, total, err := h.adminService.GetAdminUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get admin users",
		})
		return
	}
	if admins == nil {
		admins = []models.AdminUser{}
	}

	c.JSON(http.StatusOK, gin.H{
		"admins": admins,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// CreateAdminUser creates a new admin user
// @Summary Create admin user
// @Description Create a new admin user (requires super-admin privileges)
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param admin body CreateAdminUserRequest true "Admin user details"
// @Success 201 {object} models.AdminUser
// @Router /admin/users/admins [post]
func (h *AdminHandler) CreateAdminUser(c *gin.Context) {
	var req CreateAdminUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	admin := &models.AdminUser{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	createdAdmin, err := h.adminService.CreateAdminUser(admin, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdAdmin)
}

// UpdateAdminUser updates an admin user
// @Summary Update admin user
// @Description Update an existing admin user
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Admin User ID"
// @Param admin body UpdateAdminUserRequest true "Admin user updates"
// @Success 200 {object} models.AdminUser
// @Router /admin/users/admins/{id} [put]
func (h *AdminHandler) UpdateAdminUser(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid admin user ID",
		})
		return
	}

	// Prevent self-modification of role/status
	currentAdminID, _ := c.Get("user_id")
	if adminID == currentAdminID.(int) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot modify your own account",
		})
		return
	}

	var updates models.AdminUser
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = h.adminService.UpdateAdminUser(adminID, &updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return updated admin
	updatedAdmin, err := h.adminService.GetAdminUserByID(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated admin",
		})
		return
	}

	c.JSON(http.StatusOK, updatedAdmin)
}

// DeleteAdminUser deletes an admin user
// @Summary Delete admin user
// @Description Delete an admin user (cannot delete self)
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Admin User ID"
// @Success 200 {object} map[string]string
// @Router /admin/users/admins/{id} [delete]
func (h *AdminHandler) DeleteAdminUser(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid admin user ID",
		})
		return
	}

	// Prevent self-deletion
	currentAdminID, _ := c.Get("user_id")
	if adminID == currentAdminID.(int) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete your own account",
		})
		return
	}

	err = h.adminService.DeleteAdminUser(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete admin user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin user deleted successfully",
		"adminId": adminID,
	})
}

// UpdateAdminPassword changes admin password
// @Summary Change admin password
// @Description Change admin user password
// @Tags admin-auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body ChangeAdminPasswordRequest true "Password change details"
// @Success 200 {object} map[string]string
// @Router /admin/auth/password [patch]
func (h *AdminHandler) UpdateAdminPassword(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	var req ChangeAdminPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.adminService.UpdateAdminPassword(adminID.(int), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// GetTodayStats returns today's dashboard statistics
// @Summary Get today's statistics
// @Description Get dashboard statistics for today
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.AdminStats
// @Router /admin/dashboard/stats/today [get]
func (h *AdminHandler) GetTodayStats(c *gin.Context) {
	stats, err := h.adminService.GetTodayStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get today's statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Request/Response structs
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AdminLoginResponse struct {
	Admin *models.AdminUser `json:"admin"`
	Token string            `json:"token"`
}

type CreateAdminUserRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin manager staff"`
	IsActive bool   `json:"isActive"`
}

type UpdateAdminUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=admin manager staff"`
	IsActive *bool  `json:"isActive,omitempty"`
}

type ChangeAdminPasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
}
