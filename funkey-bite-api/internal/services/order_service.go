package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type OrderService interface {
	// CreateOrder validates and saves a checkout in one transaction. ctx
	// carries who placed it and its idempotency key (see middleware.RequestContext).
	CreateOrder(ctx context.Context, checkout *models.Checkout) (*models.Order, error)
	GetOrderByID(id int) (*models.Order, error)
	GetOrdersByUserID(userID int) ([]models.Order, error)
	UpdateOrderStatus(id int, status string) error
	CalculateOrderTotal(items []models.OrderItemRequest) (float64, error)
	GetOrderByOrderNumber(orderNumber string) (*models.Order, error)
	GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error)
	CancelOrder(id int, userID int) error
	MarkPaymentCollected(id int, via string) (*models.Order, error)
}

// ItemUnavailableError reports a menu item that can't be ordered in the
// requested quantity; its message is safe to show to customers.
type ItemUnavailableError struct {
	Message string
}

func (e *ItemUnavailableError) Error() string { return e.Message }

var (
	ErrOrderNotFound = errors.New("order not found")
	// ErrPromotionUnavailable is returned when a promotion ran out or ended
	// between checking it and placing the order.
	ErrPromotionUnavailable = errors.New("this promotion is no longer available")
	// ErrPaymentNotCollectable is returned for orders that are not pay on
	// delivery, or whose payment was already recorded.
	ErrPaymentNotCollectable = errors.New("only pay-on-delivery orders awaiting payment can be marked as paid")
)

type orderService struct {
	orderRepo           repository.IOrderRepository
	menuRepo            menuItemReader
	stockService        stockService
	notificationService NotificationService
	checkout            checkoutStore
}

// checkoutStore saves a checkout in a single transaction (see
// repository.CheckoutRepository).
type checkoutStore interface {
	PlaceOrder(ctx context.Context, checkout *models.Checkout) (int, error)
}

type menuItemReader interface {
	GetByID(id int) (*models.MenuItem, error)
}

// stockService gives the customer a clear message before checkout; the
// stock itself is taken, under a lock, when the order is saved.
type stockService interface {
	CheckAvailability(menuItemID, quantity int) (bool, string)
}

func NewOrderService(orderRepo repository.IOrderRepository, menuRepo menuItemReader, stockService stockService, notificationService NotificationService, checkout checkoutStore) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		menuRepo:            menuRepo,
		stockService:        stockService,
		notificationService: notificationService,
		checkout:            checkout,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, checkout *models.Checkout) (*models.Order, error) {
	names := make(map[int]string, len(checkout.Items))
	for _, item := range checkout.Items {
		menuItem, err := s.menuRepo.GetByID(item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
		}
		if menuItem == nil {
			return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
		}
		if !menuItem.IsAvailable {
			return nil, &ItemUnavailableError{Message: menuItem.Name + " is not available right now"}
		}

		if menuItem.Price != item.UnitPrice {
			return nil, fmt.Errorf("price mismatch for %s: expected ₦%.2f, got ₦%.2f",
				menuItem.Name, menuItem.Price, item.UnitPrice)
		}

		if menuItem.Name != item.Name {
			return nil, fmt.Errorf("item name mismatch for ID %d: expected '%s', got '%s'",
				item.MenuItemID, menuItem.Name, item.Name)
		}

		if ok, message := s.stockService.CheckAvailability(item.MenuItemID, item.Quantity); !ok {
			return nil, &ItemUnavailableError{Message: item.Name + " " + message}
		}
		names[item.MenuItemID] = menuItem.Name
	}

	orderID, err := s.checkout.PlaceOrder(ctx, checkout)
	var shortage *repository.InsufficientStockError
	switch {
	case errors.As(err, &shortage):
		// Stock ran out between the check above and saving the order.
		return nil, &ItemUnavailableError{Message: names[shortage.MenuItemID] + " " + shortageMessage(shortage)}
	case errors.Is(err, repository.ErrPromotionUnavailable):
		return nil, ErrPromotionUnavailable
	case err != nil:
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	createdOrder, err := s.orderRepo.GetOrderWithItems(orderID)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.notificationService.SendOrderConfirmation(createdOrder); err != nil {
			log.Printf("Failed to send order confirmation notification: %v", err)
		}
	}()

	return createdOrder, nil
}

// shortageMessage completes a sentence starting with the item's name, in the
// same words as CheckAvailability.
func shortageMessage(shortage *repository.InsufficientStockError) string {
	switch {
	case shortage.Untracked:
		return "is not available right now"
	case shortage.Available <= 0:
		return "is out of stock"
	default:
		return fmt.Sprintf("has only %d left", shortage.Available)
	}
}

// MarkPaymentCollected records that the customer paid a pay-on-delivery order.
func (s *orderService) MarkPaymentCollected(id int, via string) (*models.Order, error) {
	updated, err := s.orderRepo.MarkPaymentCollected(id, via, time.Now())
	if err != nil {
		return nil, err
	}
	if !updated {
		order, err := s.orderRepo.GetOrderWithItems(id)
		if err != nil {
			return nil, err
		}
		if order == nil {
			return nil, ErrOrderNotFound
		}
		return nil, ErrPaymentNotCollectable
	}
	return s.orderRepo.GetOrderWithItems(id)
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
