// internal/services/order_service.go
package services

import (
	"fmt"
	// "time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type OrderService interface {
	CreateOrder(order *models.Order, items []models.OrderItemRequest) (*models.Order, error)
	GetOrderByID(id int) (*models.Order, error)
	GetOrdersByUserID(userID int) ([]models.Order, error)
	UpdateOrderStatus(id int, status string) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
}

func NewOrderService(orderRepo repository.OrderRepository, menuRepo repository.MenuRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *orderService) CreateOrder(order *models.Order, items []models.OrderItemRequest) (*models.Order, error) {
	// Validate all menu items exist and are available
	for _, item := range items {
		menuItem, err := s.menuRepo.GetByID(item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
		}
		if !menuItem.IsAvailable {
			return nil, fmt.Errorf("menu item not available: %s", menuItem.Name)
		}
	}

	// Create the order
	createdOrder, err := s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItem := models.OrderItem{
			OrderID:             createdOrder.ID,
			MenuItemID:          item.MenuItemID,
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			SpecialInstructions: item.SpecialInstructions,
		}
		if _, err := s.orderRepo.CreateOrderItem(&orderItem); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
		orderItems = append(orderItems, orderItem)
	}

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
