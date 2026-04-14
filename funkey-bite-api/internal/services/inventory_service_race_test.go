package services

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type fakeInventoryStore struct {
	mu          sync.Mutex
	itemsByMenu map[int]*models.InventoryItem
}

func (f *fakeInventoryStore) GetByID(id int) (*models.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, item := range f.itemsByMenu {
		if item.ID == id {
			copy := *item
			return &copy, nil
		}
	}
	return nil, nil
}

func (f *fakeInventoryStore) GetByMenuItemID(menuItemID int) (*models.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item := f.itemsByMenu[menuItemID]
	if item == nil {
		return nil, nil
	}
	copy := *item
	return &copy, nil
}

func (f *fakeInventoryStore) GetAll() ([]models.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	items := make([]models.InventoryItem, 0, len(f.itemsByMenu))
	for _, item := range f.itemsByMenu {
		items = append(items, *item)
	}
	return items, nil
}

func (f *fakeInventoryStore) GetLowStock() ([]models.InventoryItem, error) { return nil, nil }

func (f *fakeInventoryStore) UpdateStock(itemID int, newStock int, operation string, reason string) error {
	return fmt.Errorf("UpdateStock should not be called; use AdjustStockByMenuItemID")
}

func (f *fakeInventoryStore) AdjustStockByMenuItemID(menuItemID int, quantity int, operation string, reason string) (*models.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	item := f.itemsByMenu[menuItemID]
	if item == nil {
		return nil, nil
	}

	switch operation {
	case "add":
		item.CurrentStock += quantity
	case "subtract":
		if item.CurrentStock < quantity {
			return nil, fmt.Errorf("insufficient stock. Current: %d, Requested: %d", item.CurrentStock, quantity)
		}
		item.CurrentStock -= quantity
	case "set":
		if quantity < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
		item.CurrentStock = quantity
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}
	item.UpdatedAt = time.Now()
	copy := *item
	return &copy, nil
}

func (f *fakeInventoryStore) CreateInventoryItem(item *models.InventoryItem) (*models.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *item
	f.itemsByMenu[item.MenuItemID] = &copy
	return &copy, nil
}

func (f *fakeInventoryStore) GetAlerts(resolved bool) ([]models.InventoryAlert, error) {
	return nil, nil
}
func (f *fakeInventoryStore) ResolveAlert(alertID int) error { return nil }

type fakeInventoryMenuReader struct{}

func (f *fakeInventoryMenuReader) GetByID(id int) (*models.MenuItem, error) {
	return &models.MenuItem{ID: id, Name: "Item", IsAvailable: true}, nil
}

func TestConcurrentDeductAndRestockAvoidsLostUpdates(t *testing.T) {
	store := &fakeInventoryStore{
		itemsByMenu: map[int]*models.InventoryItem{
			1: {
				ID:           10,
				MenuItemID:   1,
				Name:         "Burger",
				CurrentStock: 200,
				IsActive:     true,
			},
		},
	}

	svc := NewInventoryService(store, &fakeInventoryMenuReader{})

	const workers = 200
	var wg sync.WaitGroup
	wg.Add(workers * 2)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			if err := svc.DeductStockForOrder(1, 1, i+1); err != nil {
				t.Errorf("DeductStockForOrder() error = %v", err)
			}
		}(i)
		go func() {
			defer wg.Done()
			if err := svc.RestockItem(1, 1, "test restock"); err != nil {
				t.Errorf("RestockItem() error = %v", err)
			}
		}()
	}
	wg.Wait()

	item, err := svc.GetInventoryByMenuItemID(1)
	if err != nil {
		t.Fatalf("GetInventoryByMenuItemID() error = %v", err)
	}
	if item == nil {
		t.Fatal("inventory item not found after concurrent operations")
	}
	if item.CurrentStock != 200 {
		t.Fatalf("lost update detected: final stock = %d, want %d", item.CurrentStock, 200)
	}
}
