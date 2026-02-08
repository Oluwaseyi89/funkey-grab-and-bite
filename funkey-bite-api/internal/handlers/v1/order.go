package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type OrderHandler struct {
	orderService services.OrderService
	authService  services.AuthService
	validate     *validator.Validate
}

func NewOrderHandler(orderService services.OrderService, authService services.AuthService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		authService:  authService,
		validate:     validator.New(),
	}
}

// func (h *OrderHandler) CreateOrder(c *gin.Context) {
// 	var req models.OrderWithAuth

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	if err := h.validate.Struct(req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check authentication
// 	var userID *int

// 	// First, check if user is already authenticated via token
// 	if authUserID, exists := c.Get("user_id"); exists {
// 		id := authUserID.(int)
// 		userID = &id
// 	} else {
// 		// Authenticate via order data (phone + password)
// 		user, err := h.authService.AuthenticateOrder(req)
// 		if err != nil {
// 			switch err.Error() {
// 			case "user_not_found":
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "User not found. Please create an account first.",
// 					"code":  "USER_NOT_FOUND",
// 				})
// 				return
// 			case "password_required":
// 				c.JSON(http.StatusUnauthorized, gin.H{
// 					"error": "Password required for existing user",
// 					"code":  "PASSWORD_REQUIRED",
// 				})
// 				return
// 			case "invalid_password":
// 				c.JSON(http.StatusUnauthorized, gin.H{
// 					"error": "Invalid password",
// 					"code":  "INVALID_PASSWORD",
// 				})
// 				return
// 			default:
// 				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 				return
// 			}
// 		}
// 		userID = &user.ID
// 	}

// 	// Create order using the updated model
// 	order := &models.Order{
// 		OrderNumber:   utils.GenerateOrderNumber(),
// 		UserID:        userID,
// 		CustomerName:  req.CustomerName,
// 		CustomerPhone: req.CustomerPhone,
// 		CustomerEmail: req.CustomerEmail,
// 		OrderType:     req.OrderType,
// 		Status:        models.OrderStatusPending,
// 		TotalAmount:   calculateTotal(req.Items),
// 		Notes:         req.Notes,
// 		PickupTime:    req.PickupTime,
// 		CreatedAt:     time.Now(),
// 	}

// 	createdOrder, err := h.orderService.CreateOrder(order, req.Items)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
// 		return
// 	}

// 	// Return the order directly (it will use the json tags from the model)
// 	c.JSON(http.StatusCreated, createdOrder)
// }

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.OrderWithAuth

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to calculate total with price validation
	totalAmount, err := h.orderService.CalculateOrderTotal(req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check authentication
	var userID *int

	// First, check if user is already authenticated via token
	if authUserID, exists := c.Get("user_id"); exists {
		id := authUserID.(int)
		userID = &id
	} else {
		// Authenticate via order data (phone + password)
		user, err := h.authService.AuthenticateOrder(req)
		if err != nil {
			switch err.Error() {
			case "user_not_found":
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "User not found. Please create an account first.",
					"code":  "USER_NOT_FOUND",
				})
				return
			case "password_required":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Password required for existing user",
					"code":  "PASSWORD_REQUIRED",
				})
				return
			case "invalid_password":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid password",
					"code":  "INVALID_PASSWORD",
				})
				return
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		userID = &user.ID
	}

	// Create order using the validated total
	order := &models.Order{
		OrderNumber:   utils.GenerateOrderNumber(),
		UserID:        userID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		OrderType:     req.OrderType,
		Status:        models.OrderStatusPending,
		TotalAmount:   totalAmount, // Use validated total from service
		Notes:         req.Notes,
		PickupTime:    req.PickupTime,
		CreatedAt:     time.Now(),
	}

	createdOrder, err := h.orderService.CreateOrder(order, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	// Return the order directly (it will use the json tags from the model)
	c.JSON(http.StatusCreated, createdOrder)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	// TODO: Implement order retrieval
	c.JSON(http.StatusOK, gin.H{"message": "Get order " + orderID})
}

// GetUserOrders gets all orders for the authenticated user
// @Summary Get user orders
// @Description Get all orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /auth/orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	// Get orders from service
	orders, err := h.orderService.GetOrdersByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user orders",
		})
		return
	}

	// Simple pagination
	total := len(orders)
	offset := (page - 1) * limit
	if offset >= total {
		c.JSON(http.StatusOK, gin.H{
			"orders": []models.Order{},
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": (total + limit - 1) / limit,
			},
		})
		return
	}

	end := offset + limit
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders[offset:end],
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// GetUserOrder gets a specific order for the authenticated user
// @Summary Get user order by ID
// @Description Get a specific order for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Router /auth/orders/{id} [get]
func (h *OrderHandler) GetUserOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Get order
	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	// Check if order belongs to user
	if order.UserID == nil || *order.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) TrackOrder(c *gin.Context) {
	orderNumber := c.Param("orderNumber")

	order, err := h.orderService.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	// Check if user is authorized to view this order
	userID, exists := c.Get("user_id")
	if exists {
		// Authenticated user - can view their own orders
		if order.UserID != nil && *order.UserID == userID.(int) {
			c.JSON(http.StatusOK, order)
			return
		}
	}

	// For non-owners, return limited information
	c.JSON(http.StatusOK, gin.H{
		"orderNumber":        order.OrderNumber,
		"status":             order.Status,
		"estimatedReadyTime": order.EstimatedReadyTime,
		"createdAt":          order.CreatedAt,
		"message":            "Limited order information available",
	})
}

func (h *OrderHandler) TrackOrderPublic(c *gin.Context) {
	phone := c.Param("phone")
	orderNumber := c.Param("orderNumber")

	// Basic validation
	if phone == "" || orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone and order number are required",
		})
		return
	}

	// Rate limiting check (simple version - can enhance later)
	// Consider adding actual rate limiting middleware

	order, err := h.orderService.GetOrderByPhoneAndOrderNumber(phone, orderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	if order == nil {
		// Don't reveal if order exists but credentials are wrong
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found or phone number doesn't match",
		})
		return
	}

	// Return public-safe order information (no sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"orderNumber":        order.OrderNumber,
		"status":             order.Status,
		"estimatedReadyTime": order.EstimatedReadyTime,
		"pickupTime":         order.PickupTime,
		"createdAt":          order.CreatedAt,
		"items":              order.Items,
		"orderType":          order.OrderType,
		"totalAmount":        order.TotalAmount,
		// Omit sensitive fields: customerName, customerPhone, customerEmail, notes
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req struct {
		Status string `json:"status" validate:"required,oneof=pending confirmed preparing ready completed cancelled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// TODO: Implement status update (admin only)
	c.JSON(http.StatusOK, gin.H{
		"orderId": orderID,
		"status":  req.Status,
		"message": "Order status updated",
	})
}

// CancelOrder cancels a pending order
// @Summary Cancel order
// @Description Cancel a pending order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string
// @Router /orders/{id}/cancel [patch]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderService.CancelOrder(orderID, userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"orderId": orderID,
		"status":  "cancelled",
	})
}

// func calculateTotal(items []models.OrderItemRequest) float64 {
// 	total := 0.0
// 	for _, item := range items {
// 		total += item.UnitPrice * float64(item.Quantity)
// 	}
// 	return total
// }

// // Helper method for order conversion if needed
// func (h *OrderHandler) convertOrderForFrontend(order *models.Order) map[string]interface{} {
// 	if order == nil {
// 		return nil
// 	}

// 	return map[string]interface{}{
// 		"id":            order.ID,
// 		"orderNumber":   order.OrderNumber,
// 		"customerName":  order.CustomerName,
// 		"customerPhone": order.CustomerPhone,
// 		"customerEmail": order.CustomerEmail,
// 		"orderType":     string(order.OrderType),
// 		"status":        string(order.Status),
// 		"totalAmount":   order.TotalAmount,
// 		"notes":         order.Notes,
// 		"pickupTime":    order.PickupTime,
// 		"createdAt":     order.CreatedAt,
// 		"items":         order.Items,
// 	}
// }
