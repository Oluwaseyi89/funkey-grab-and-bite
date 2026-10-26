package services

import (
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

func TestCheckAvailabilityMessages(t *testing.T) {
	store := &fakeInventoryStore{itemsByMenu: map[int]*models.InventoryItem{
		1: {ID: 1, MenuItemID: 1, CurrentStock: 0, IsActive: true},
		2: {ID: 2, MenuItemID: 2, CurrentStock: 2, IsActive: true},
		3: {ID: 3, MenuItemID: 3, CurrentStock: 9, IsActive: false},
	}}
	svc := NewInventoryService(store, &fakeInventoryMenuReader{})

	cases := []struct {
		name       string
		menuItemID int
		quantity   int
		wantOK     bool
		wantMsg    string
	}{
		{"out of stock", 1, 1, false, "is out of stock"},
		{"not enough left", 2, 3, false, "has only 2 left"},
		{"enough stock", 2, 2, true, ""},
		{"deactivated", 3, 1, false, "is not available right now"},
		{"no stock record", 4, 1, false, "is not available right now"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, msg := svc.CheckAvailability(tc.menuItemID, tc.quantity)
			if ok != tc.wantOK || msg != tc.wantMsg {
				t.Fatalf("CheckAvailability() = %t, %q; want %t, %q", ok, msg, tc.wantOK, tc.wantMsg)
			}
		})
	}
}

func TestUpdateStockAllowsSettingZeroButNotAddingZero(t *testing.T) {
	store := &fakeInventoryStore{itemsByMenu: map[int]*models.InventoryItem{
		1: {ID: 1, MenuItemID: 1, CurrentStock: 7, IsActive: true},
	}}
	svc := NewInventoryService(store, &fakeInventoryMenuReader{})

	item, err := svc.UpdateStock(&models.InventoryUpdate{MenuItemID: 1, Quantity: 0, Operation: "set", Reason: "sold out"})
	if err != nil || item.CurrentStock != 0 {
		t.Fatalf("setting stock to 0 should succeed, got %+v, %v", item, err)
	}

	for _, operation := range []string{"add", "subtract"} {
		if _, err := svc.UpdateStock(&models.InventoryUpdate{MenuItemID: 1, Quantity: 0, Operation: operation, Reason: "noop"}); err == nil {
			t.Fatalf("%s with quantity 0 should be rejected", operation)
		}
	}
}
