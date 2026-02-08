package services

import (
	"fmt"
	"log"

	// "time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type OrderService interface {
	CreateOrder(order *models.Order, items []models.OrderItemRequest) (*models.Order, error)
	GetOrderByID(id int) (*models.Order, error)
	GetOrdersByUserID(userID int) ([]models.Order, error)
	UpdateOrderStatus(id int, status string) error
	CalculateOrderTotal(items []models.OrderItemRequest) (float64, error)
	GetOrderByOrderNumber(orderNumber string) (*models.Order, error)
	GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error)
	CancelOrder(id int, userID int) error
}

type orderService struct {
	// orderRepo repository.OrderRepository
	orderRepo           repository.IOrderRepository // Change to interface
	menuRepo            repository.MenuRepository
	notificationService NotificationService // Add this

}

func NewOrderService(orderRepo repository.IOrderRepository, menuRepo repository.MenuRepository, notificationService NotificationService) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		menuRepo:            menuRepo,
		notificationService: notificationService,
	}
}

func (s *orderService) CreateOrder(order *models.Order, items []models.OrderItemRequest) (*models.Order, error) {
	// Validate all menu items exist, are available, and prices match
	for _, item := range items {
		menuItem, err := s.menuRepo.GetByID(item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
		}
		if menuItem == nil {
			return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
		}
		if !menuItem.IsAvailable {
			return nil, fmt.Errorf("menu item not available: %s", menuItem.Name)
		}

		// CRITICAL: Verify price matches current menu price
		if menuItem.Price != item.UnitPrice {
			return nil, fmt.Errorf("price mismatch for %s: expected $%.2f, got $%.2f",
				menuItem.Name, menuItem.Price, item.UnitPrice)
		}

		// Also verify the name matches (extra security)
		if menuItem.Name != item.Name {
			return nil, fmt.Errorf("item name mismatch for ID %d: expected '%s', got '%s'",
				item.MenuItemID, menuItem.Name, item.Name)
		}
	}

	// Start database transaction
	tx, err := s.orderRepo.BeginTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Rollback in case of error
	var txErr error
	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	// Create order within transaction
	createdOrder, err := s.orderRepo.CreateOrderWithTransaction(tx, order)
	if err != nil {
		txErr = err
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items within transaction
	for _, item := range items {
		orderItem := models.OrderItem{
			OrderID:             createdOrder.ID,
			MenuItemID:          item.MenuItemID,
			Name:                item.Name,
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			SpecialInstructions: item.SpecialInstructions,
		}
		_, err := s.orderRepo.CreateOrderItemWithTransaction(tx, &orderItem)
		if err != nil {
			txErr = err
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func() {
		if err := s.notificationService.SendOrderConfirmation(createdOrder); err != nil {
			log.Printf("Failed to send order confirmation notification: %v", err)
		}
	}()

	// Fetch complete order with items
	return s.orderRepo.GetOrderWithItems(createdOrder.ID)
}

func (s *orderService) GetOrderByID(id int) (*models.Order, error) {
	return s.orderRepo.GetOrderWithItems(id)
}

func (s *orderService) GetOrdersByUserID(userID int) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByUserID(userID)
}

func (s *orderService) UpdateOrderStatus(id int, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"preparing": true,
		"ready":     true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.orderRepo.UpdateOrderStatus(id, status)
}

func (s *orderService) CalculateOrderTotal(items []models.OrderItemRequest) (float64, error) {
	total := 0.0
	for _, item := range items {
		// Verify each item exists and get price
		menuItem, err := s.menuRepo.GetByID(item.MenuItemID)
		if err != nil || menuItem == nil {
			return 0, fmt.Errorf("invalid menu item ID: %d", item.MenuItemID)
		}
		if !menuItem.IsAvailable {
			return 0, fmt.Errorf("menu item not available: %s", menuItem.Name)
		}
		total += menuItem.Price * float64(item.Quantity)
	}
	return total, nil
}

func (s *orderService) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by order number: %w", err)
	}
	return order, nil
}

func (s *orderService) GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderByPhoneAndOrderNumber(phone, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by phone and order number: %w", err)
	}
	return order, nil
}

func (s *orderService) CancelOrder(id int, userID int) error {
	// First get the order to verify ownership
	order, err := s.orderRepo.GetOrderWithItems(id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	// Check if order belongs to user
	if order.UserID == nil || *order.UserID != userID {
		return fmt.Errorf("unauthorized to cancel this order")
	}

	// Check if order can be cancelled (only pending orders)
	if order.Status != models.OrderStatusPending {
		return fmt.Errorf("order cannot be cancelled (status: %s)", order.Status)
	}

	// Cancel the order
	return s.orderRepo.CancelOrder(id)
}
