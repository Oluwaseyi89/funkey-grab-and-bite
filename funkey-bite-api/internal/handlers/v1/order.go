package v1

import (
	"net/http"
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

	// Create order using the updated model
	order := &models.Order{
		OrderNumber:   utils.GenerateOrderNumber(),
		UserID:        userID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		OrderType:     req.OrderType,
		Status:        models.OrderStatusPending,
		TotalAmount:   calculateTotal(req.Items),
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

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// TODO: Implement user orders retrieval
	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"orders": []interface{}{},
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

func calculateTotal(items []models.OrderItemRequest) float64 {
	total := 0.0
	for _, item := range items {
		total += item.UnitPrice * float64(item.Quantity)
	}
	return total
}

// Helper method for order conversion if needed
func (h *OrderHandler) convertOrderForFrontend(order *models.Order) map[string]interface{} {
	if order == nil {
		return nil
	}

	return map[string]interface{}{
		"id":            order.ID,
		"orderNumber":   order.OrderNumber,
		"customerName":  order.CustomerName,
		"customerPhone": order.CustomerPhone,
		"customerEmail": order.CustomerEmail,
		"orderType":     string(order.OrderType),
		"status":        string(order.Status),
		"totalAmount":   order.TotalAmount,
		"notes":         order.Notes,
		"pickupTime":    order.PickupTime,
		"createdAt":     order.CreatedAt,
		"items":         order.Items,
	}
}
