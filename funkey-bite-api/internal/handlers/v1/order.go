package v1

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/realtime"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type OrderHandler struct {
	orderService services.OrderService
	authService  services.AuthService
	settings     services.SettingsService
	promotions   services.PromotionService
	validate     *validator.Validate
}

func NewOrderHandler(orderService services.OrderService, authService services.AuthService, settingsService services.SettingsService, promotionService services.PromotionService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		authService:  authService,
		settings:     settingsService,
		promotions:   promotionService,
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

	totalAmount, err := h.orderService.CalculateOrderTotal(req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok, message := h.settings.ValidateMinimumOrder(totalAmount, string(req.OrderType)); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	if ok, message := h.settings.ValidateOrderTime(string(req.OrderType), req.PickupTime); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	if ok, message := h.settings.CanAcceptOrders(string(req.OrderType)); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	_, _, totalWithFees := h.settings.CalculateOrderFees(totalAmount, string(req.OrderType))

	var userID *int

	if authUserID, exists := c.Get("user_id"); exists {
		id := authUserID.(int)
		userID = &id
	} else {
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

	finalTotal := totalWithFees
	promotionCode := ""
	if req.PromotionCode != nil {
		promotionCode = strings.TrimSpace(*req.PromotionCode)
	}

	if promotionCode != "" {
		validation, err := h.promotions.ValidatePromotion(promotionCode, totalWithFees, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate promotion"})
			return
		}

		if !validation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": validation.Message})
			return
		}

		finalTotal = totalWithFees - validation.Discount
		if finalTotal < 0 {
			finalTotal = 0
		}
	}

	order := &models.Order{
		OrderNumber:   utils.GenerateOrderNumber(),
		UserID:        userID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		OrderType:     req.OrderType,
		Status:        models.OrderStatusPending,
		TotalAmount:   finalTotal,
		Notes:         req.Notes,
		PickupTime:    req.PickupTime,
		CreatedAt:     time.Now(),
	}

	createdOrder, err := h.orderService.CreateOrder(order, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	if promotionCode != "" {
		if _, err := h.promotions.ApplyPromotionByCode(promotionCode, createdOrder.ID, userID, totalWithFees); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply promotion: " + err.Error()})
			return
		}
	}

	realtime.GlobalHub.Broadcast("new_order", createdOrder)

	c.JSON(http.StatusCreated, createdOrder)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
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

	orders, err := h.orderService.GetOrdersByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user orders",
		})
		return
	}

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

	userID, exists := c.Get("user_id")
	if exists {
		if order.UserID != nil && *order.UserID == userID.(int) {
			c.JSON(http.StatusOK, order)
			return
		}
	}

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

	if phone == "" || orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone and order number are required",
		})
		return
	}

	order, err := h.orderService.GetOrderByPhoneAndOrderNumber(phone, orderNumber)
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

	c.JSON(http.StatusOK, gin.H{
		"orderNumber":        order.OrderNumber,
		"status":             order.Status,
		"estimatedReadyTime": order.EstimatedReadyTime,
		"pickupTime":         order.PickupTime,
		"createdAt":          order.CreatedAt,
		"orderType":          order.OrderType,
		"message":            "Limited order information available",
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
