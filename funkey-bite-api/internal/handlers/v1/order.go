package v1

import (
	"errors"
	"log"
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
	orderService   services.OrderService
	authService    services.AuthService
	settings       services.SettingsService
	promotions     services.PromotionService
	paymentService services.PaymentService
	invoices       services.InvoiceService
	validate       *validator.Validate
}

func NewOrderHandler(orderService services.OrderService, authService services.AuthService, settingsService services.SettingsService, promotionService services.PromotionService, paymentService services.PaymentService, invoiceService services.InvoiceService) *OrderHandler {
	return &OrderHandler{
		orderService:   orderService,
		authService:    authService,
		settings:       settingsService,
		promotions:     promotionService,
		paymentService: paymentService,
		invoices:       invoiceService,
		validate:       validator.New(),
	}
}

// CreateOrder places an order. Customers check out as guests: name, phone and
// email are required, plus an address for delivery. Pay on delivery is only
// accepted while the admin has it enabled, and raises an invoice.
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.OrderWithAuth

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerPhone = utils.CleanPhone(req.CustomerPhone)
	if req.CustomerEmail != nil {
		email := strings.TrimSpace(*req.CustomerEmail)
		req.CustomerEmail = &email
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidPhone(req.CustomerPhone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid phone number"})
		return
	}

	var deliveryAddress *string
	if req.OrderType == models.OrderTypeDelivery {
		if req.DeliveryAddress == nil || len(strings.TrimSpace(*req.DeliveryAddress)) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A delivery address is required for delivery orders"})
			return
		}
		address := strings.TrimSpace(*req.DeliveryAddress)
		deliveryAddress = &address
	}

	if req.PaymentMethod == models.PaymentMethodCash {
		settings, err := h.settings.GetSettings()
		if err != nil || settings == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to check payment options"})
			return
		}
		if !settings.PayOnDeliveryEnabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pay on delivery is not available right now. Please pay online."})
			return
		}
	}

	subtotal, err := h.orderService.CalculateOrderTotal(req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok, message := h.settings.ValidateMinimumOrder(subtotal, string(req.OrderType)); !ok {
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

	deliveryFee, taxAmount, totalWithFees := h.settings.CalculateOrderFees(subtotal, string(req.OrderType))

	// A returning guest is linked now so per-customer promotion limits apply.
	// A first-time guest's record is only created once the order is placed,
	// so a rejected checkout leaves no customer behind.
	var userID *int
	if authUserID, exists := c.Get("user_id"); exists {
		id := authUserID.(int)
		userID = &id
	} else if customer, err := h.authService.FindOrderCustomer(req.CustomerPhone); err != nil {
		log.Printf("Failed to look up customer for phone ending %s: %v", lastDigits(req.CustomerPhone), err)
	} else if customer != nil {
		userID = &customer.ID
	}

	finalTotal := totalWithFees
	discount := 0.0
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

		discount = validation.Discount
		finalTotal = totalWithFees - discount
		if finalTotal < 0 {
			finalTotal = 0
		}
	}

	orderNumber := utils.GenerateOrderNumber()
	var invoiceNumber *string
	if req.PaymentMethod == models.PaymentMethodCash {
		number := services.InvoiceNumberFor(orderNumber)
		invoiceNumber = &number
	}

	order := &models.Order{
		OrderNumber:   orderNumber,
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
		PaymentMethod: req.PaymentMethod,
		// Both methods start pending: online payments until Paystack confirms,
		// pay on delivery until the admin records the payment.
		PaymentStatus:   models.PaymentStatusPending,
		DeliveryAddress: deliveryAddress,
		InvoiceNumber:   invoiceNumber,
		SubtotalAmount:  &subtotal,
		DeliveryFee:     &deliveryFee,
		TaxAmount:       &taxAmount,
		DiscountAmount:  &discount,
	}

	createdOrder, err := h.orderService.CreateOrder(order, req.Items)
	var unavailable *services.ItemUnavailableError
	if errors.As(err, &unavailable) {
		c.JSON(http.StatusConflict, gin.H{"error": unavailable.Message})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	if userID == nil {
		userID = h.recordGuestCustomer(createdOrder, req)
	}

	if promotionCode != "" {
		if _, err := h.promotions.ApplyPromotionByCode(promotionCode, createdOrder.ID, userID, totalWithFees); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply promotion: " + err.Error()})
			return
		}
	}

	if createdOrder.PaymentMethod == models.PaymentMethodTransfer {
		// The order already exists at this point (items + stock committed); a
		// Paystack outage here shouldn't lose the order, just leave payment
		// pending with no account details for the customer to retry against later.
		if withPayment, err := h.paymentService.InitiateTransferPayment(createdOrder); err != nil {
			log.Printf("Failed to initiate transfer payment for order #%s: %v", createdOrder.OrderNumber, err)
		} else {
			createdOrder = withPayment
		}
	}

	if createdOrder.InvoiceNumber != nil {
		invoiceOrder := createdOrder
		go func() {
			if err := h.invoices.SendInvoice(invoiceOrder); err != nil {
				log.Printf("Failed to send invoice for order #%s: %v", invoiceOrder.OrderNumber, err)
			}
		}()
	}

	realtime.GlobalBroadcaster.Broadcast("new_order", createdOrder)

	c.JSON(http.StatusCreated, createdOrder)
}

// recordGuestCustomer creates the customer record for a first-time guest and
// links the placed order to it. The record is for the admin's follow-up;
// failing to save it must never lose the order itself.
func (h *OrderHandler) recordGuestCustomer(order *models.Order, req models.OrderWithAuth) *int {
	customer, err := h.authService.ResolveOrderCustomer(req.CustomerName, req.CustomerPhone, req.CustomerEmail)
	if err != nil || customer == nil {
		log.Printf("Failed to record customer for phone ending %s: %v", lastDigits(req.CustomerPhone), err)
		return nil
	}
	if err := h.orderService.AssignCustomer(order.ID, customer.ID); err != nil {
		log.Printf("Failed to link order #%s to its customer: %v", order.OrderNumber, err)
		return nil
	}
	order.UserID = &customer.ID
	return &customer.ID
}

// lastDigits keeps phone numbers out of logs apart from a short suffix.
func lastDigits(phone string) string {
	if len(phone) <= 4 {
		return phone
	}
	return phone[len(phone)-4:]
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
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

// LookupOrder lets a customer see their order and its invoice using the order
// number and the phone number it was placed with. Both are sent in the body so
// the phone number never appears in URLs or access logs.
// @Summary Look up an order
// @Description Get an order and its invoice by order number and customer phone
// @Tags orders
// @Accept json
// @Produce json
// @Param lookup body models.OrderLookupRequest true "Order number and phone"
// @Success 200 {object} map[string]interface{}
// @Router /orders/lookup [post]
func (h *OrderHandler) LookupOrder(c *gin.Context) {
	var req models.OrderLookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order number and phone are required"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order number and phone are required"})
		return
	}

	order, err := h.orderService.GetOrderByOrderNumber(strings.TrimSpace(req.OrderNumber))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}
	// The same answer for an unknown order and a wrong phone, so the lookup
	// can't be used to confirm which order numbers exist.
	if order == nil || !utils.SamePhone(order.CustomerPhone, req.Phone) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	invoice, err := h.invoices.BuildInvoice(order)
	if err != nil {
		log.Printf("Failed to build invoice for order #%s: %v", order.OrderNumber, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"order":   order,
		"invoice": invoice,
	})
}

// GetOrderInvoice returns the invoice of an order for the admin
// @Summary Get order invoice
// @Description Get the invoice raised for a pay-on-delivery order
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.Invoice
// @Router /admin/orders/{id}/invoice [get]
func (h *OrderHandler) GetOrderInvoice(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	invoice, err := h.invoices.BuildInvoice(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build invoice"})
		return
	}
	if invoice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This order has no invoice"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// MarkPaymentCollected records that a pay-on-delivery order was paid
// @Summary Mark pay-on-delivery order as paid
// @Description Record that the customer paid on delivery or at pickup, and how
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param payment body models.PaymentCollectionRequest true "How the payment was collected"
// @Success 200 {object} models.Order
// @Router /admin/orders/{id}/payment [patch]
func (h *OrderHandler) MarkPaymentCollected(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req models.PaymentCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectedVia must be cash, pos or transfer"})
		return
	}

	order, err := h.orderService.MarkPaymentCollected(orderID, req.CollectedVia)
	switch {
	case errors.Is(err, services.ErrOrderNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	case errors.Is(err, services.ErrPaymentNotCollectable):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment"})
		return
	}

	realtime.GlobalBroadcaster.Broadcast("order_updated", order)

	c.JSON(http.StatusOK, order)
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
