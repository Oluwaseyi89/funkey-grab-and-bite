package services

import (
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type InventoryService interface {
	GetInventoryItem(id int) (*models.InventoryItem, error)
	GetInventoryByMenuItemID(menuItemID int) (*models.InventoryItem, error)
	GetAllInventory() ([]models.InventoryItem, error)
	GetLowStock() ([]models.InventoryItem, error)
	UpdateStock(update *models.InventoryUpdate) (*models.InventoryItem, error)
	CreateInventoryItem(item *models.InventoryItem) (*models.InventoryItem, error)
	GetAlerts(resolved bool) ([]models.InventoryAlert, error)
	ResolveAlert(alertID int) error
	CheckAvailability(menuItemID, quantity int) (bool, string)
	DeductStockForOrder(menuItemID, quantity int, orderID int) error
	RestockItem(menuItemID, quantity int, reason string) error
	GetInventoryDashboard() (map[string]interface{}, error)
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	menuRepo      repository.MenuRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository, menuRepo repository.MenuRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		menuRepo:      menuRepo,
	}
}

func (s *inventoryService) GetInventoryItem(id int) (*models.InventoryItem, error) {
	return s.inventoryRepo.GetByID(id)
}

func (s *inventoryService) GetInventoryByMenuItemID(menuItemID int) (*models.InventoryItem, error) {
	return s.inventoryRepo.GetByMenuItemID(menuItemID)
}

func (s *inventoryService) GetAllInventory() ([]models.InventoryItem, error) {
	return s.inventoryRepo.GetAll()
}

func (s *inventoryService) GetLowStock() ([]models.InventoryItem, error) {
	return s.inventoryRepo.GetLowStock()
}

func (s *inventoryService) UpdateStock(update *models.InventoryUpdate) (*models.InventoryItem, error) {
	item, err := s.inventoryRepo.GetByMenuItemID(update.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}
	if item == nil {
		return nil, fmt.Errorf("inventory item not found for menu item ID: %d", update.MenuItemID)
	}

	var newStock int
	switch update.Operation {
	case "add":
		newStock = item.CurrentStock + update.Quantity
	case "subtract":
		newStock = item.CurrentStock - update.Quantity
		if newStock < 0 {
			return nil, fmt.Errorf("insufficient stock. Current: %d, Requested: %d", item.CurrentStock, update.Quantity)
		}
	case "set":
		newStock = update.Quantity
		if newStock < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
	default:
		return nil, fmt.Errorf("invalid operation: %s", update.Operation)
	}

	err = s.inventoryRepo.UpdateStock(item.ID, newStock, update.Operation, update.Reason)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	return s.inventoryRepo.GetByID(item.ID)
}

func (s *inventoryService) CreateInventoryItem(item *models.InventoryItem) (*models.InventoryItem, error) {
	menuItem, err := s.menuRepo.GetByID(item.MenuItemID)
	if err != nil || menuItem == nil {
		return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
	}

	existing, err := s.inventoryRepo.GetByMenuItemID(item.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing inventory: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("inventory already exists for menu item ID: %d", item.MenuItemID)
	}

	if item.Name == "" {
		item.Name = menuItem.Name
	}
	if item.Unit == "" {
		item.Unit = "pieces"
	}
	if item.MinimumStock == 0 {
		item.MinimumStock = 10
	}
	if item.ReorderPoint == 0 {
		item.ReorderPoint = 5
	}

	return s.inventoryRepo.CreateInventoryItem(item)
}

func (s *inventoryService) GetAlerts(resolved bool) ([]models.InventoryAlert, error) {
	return s.inventoryRepo.GetAlerts(resolved)
}

func (s *inventoryService) ResolveAlert(alertID int) error {
	return s.inventoryRepo.ResolveAlert(alertID)
}

func (s *inventoryService) CheckAvailability(menuItemID, quantity int) (bool, string) {
	item, err := s.inventoryRepo.GetByMenuItemID(menuItemID)
	if err != nil || item == nil {
		return false, "Item not found in inventory"
	}

	if !item.IsActive {
		return false, "Item is not available"
	}

	if item.CurrentStock < quantity {
		return false, fmt.Sprintf("Insufficient stock. Available: %d, Requested: %d", item.CurrentStock, quantity)
	}

	return true, ""
}

func (s *inventoryService) DeductStockForOrder(menuItemID, quantity int, orderID int) error {
	update := &models.InventoryUpdate{
		MenuItemID: menuItemID,
		Quantity:   quantity,
		Operation:  "subtract",
		Reason:     fmt.Sprintf("Order #%d", orderID),
	}

	_, err := s.UpdateStock(update)
	return err
}

func (s *inventoryService) RestockItem(menuItemID, quantity int, reason string) error {
	update := &models.InventoryUpdate{
		MenuItemID: menuItemID,
		Quantity:   quantity,
		Operation:  "add",
		Reason:     reason,
	}

	_, err := s.UpdateStock(update)
	return err
}

// GetInventoryDashboard gets dashboard statistics
func (s *inventoryService) GetInventoryDashboard() (map[string]interface{}, error) {
	allItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	lowStock, err := s.inventoryRepo.GetLowStock()
	if err != nil {
		return nil, err
	}

	activeAlerts, err := s.inventoryRepo.GetAlerts(false)
	if err != nil {
		return nil, err
	}

	totalValue := 0.0
	for _, item := range allItems {
		totalValue += float64(item.CurrentStock)
	}

	return map[string]interface{}{
		"totalItems":       len(allItems),
		"lowStockItems":    len(lowStock),
		"activeAlerts":     len(activeAlerts),
		"outOfStock":       len(lowStock) - len(activeAlerts),
		"totalStockValue":  totalValue,
		"lowStockList":     lowStock,
		"activeAlertsList": activeAlerts,
	}, nil
}
