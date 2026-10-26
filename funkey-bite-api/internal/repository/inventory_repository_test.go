package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var inventoryColumns = []string{
	"id", "menu_item_id", "name", "current_stock", "minimum_stock", "reorder_point",
	"unit", "is_active", "last_restocked", "created_at", "updated_at",
}

// Stock records created alongside a menu item have never been restocked; a
// NULL last_restocked used to fail the scan and empty the whole inventory list.
func TestInventoryReadsNeverRestockedItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewInventoryRepository(db)
	now := time.Now()

	mock.ExpectQuery("FROM inventory_items").WillReturnRows(sqlmock.NewRows(inventoryColumns).
		AddRow(1, 10, "Shawarma", 0, 10, 5, "pieces", true, nil, now, now).
		AddRow(2, 11, "Chips", 25, 10, 5, "pieces", true, now, now, now))

	items, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(items) != 2 || items[0].LastRestocked != nil || items[1].LastRestocked == nil {
		t.Fatalf("unexpected items: %+v", items)
	}

	mock.ExpectQuery("FROM inventory_items").WithArgs(10).WillReturnRows(sqlmock.NewRows(inventoryColumns).
		AddRow(1, 10, "Shawarma", 0, 10, 5, "pieces", true, nil, now, now))

	item, err := repo.GetByMenuItemID(10)
	if err != nil || item == nil || item.LastRestocked != nil {
		t.Fatalf("GetByMenuItemID() = %+v, %v", item, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
