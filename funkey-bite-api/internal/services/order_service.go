package services

import (
	"fmt"
	"log"

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
	orderRepo           repository.IOrderRepository
	menuRepo            repository.MenuRepository
	notificationService NotificationService
}

func NewOrderService(orderRepo repository.IOrderRepository, menuRepo repository.MenuRepository, notificationService NotificationService) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		menuRepo:            menuRepo,
		notificationService: notificationService,
	}
}

func (s *orderService) CreateOrder(order *models.Order, items []models.OrderItemRequest) (*models.Order, error) {
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

		if menuItem.Price != item.UnitPrice {
			return nil, fmt.Errorf("price mismatch for %s: expected $%.2f, got $%.2f",
				menuItem.Name, menuItem.Price, item.UnitPrice)
		}

		if menuItem.Name != item.Name {
			return nil, fmt.Errorf("item name mismatch for ID %d: expected '%s', got '%s'",
				item.MenuItemID, menuItem.Name, item.Name)
		}
	}

	tx, err := s.orderRepo.BeginTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var txErr error
	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	createdOrder, err := s.orderRepo.CreateOrderWithTransaction(tx, order)
	if err != nil {
		txErr = err
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func() {
		if err := s.notificationService.SendOrderConfirmation(createdOrder); err != nil {
			log.Printf("Failed to send order confirmation notification: %v", err)
		}
	}()

	return s.orderRepo.GetOrderWithItems(createdOrder.ID)
}

func (s *orderService) GetOrderByID(id int) (*models.Order, error) {
	return s.orderRepo.GetOrderWithItems(id)
}

func (s *orderService) GetOrdersByUserID(userID int) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByUserID(userID)
}

func (s *orderService) UpdateOrderStatus(id int, status string) error {
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
	order, err := s.orderRepo.GetOrderWithItems(id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	if order.UserID == nil || *order.UserID != userID {
		return fmt.Errorf("unauthorized to cancel this order")
	}

	if order.Status != models.OrderStatusPending {
		return fmt.Errorf("order cannot be cancelled (status: %s)", order.Status)
	}

	return s.orderRepo.CancelOrder(id)
}
