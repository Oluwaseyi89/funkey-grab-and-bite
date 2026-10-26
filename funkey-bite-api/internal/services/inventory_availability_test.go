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
