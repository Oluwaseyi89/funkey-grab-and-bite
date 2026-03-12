package models

import (
	"time"
)

type InventoryItem struct {
	ID            int       `json:"id" db:"id"`
	MenuItemID    int       `json:"menuItemId" db:"menu_item_id"`
	Name          string    `json:"name" db:"name"`
	CurrentStock  int       `json:"currentStock" db:"current_stock"`
	MinimumStock  int       `json:"minimumStock" db:"minimum_stock"`
	ReorderPoint  int       `json:"reorderPoint" db:"reorder_point"`
	Unit          string    `json:"unit" db:"unit"` 
	IsActive      bool      `json:"isActive" db:"is_active"`
	LastRestocked time.Time `json:"lastRestocked" db:"last_restocked"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type InventoryUpdate struct {
	MenuItemID int    `json:"menuItemId" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required"`
	Operation  string `json:"operation" validate:"required,oneof=add subtract set"`
	Reason     string `json:"reason" validate:"required"`
	Notes      string `json:"notes,omitempty"`
}

type InventoryAlert struct {
	ID              int        `json:"id" db:"id"`
	InventoryItemID int        `json:"inventoryItemId" db:"inventory_item_id"`
	AlertType       string     `json:"alertType" db:"alert_type"`
	Message         string     `json:"message" db:"message"`
	IsResolved      bool       `json:"isResolved" db:"is_resolved"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	ResolvedAt      *time.Time `json:"resolvedAt,omitempty" db:"resolved_at"`
}

type InventoryHistory struct {
	ID              int       `json:"id" db:"id"`
	InventoryItemID int       `json:"inventoryItemId" db:"inventory_item_id"`
	PreviousStock   int       `json:"previousStock" db:"previous_stock"`
	NewStock        int       `json:"newStock" db:"new_stock"`
	Change          int       `json:"change" db:"change"`
	Operation       string    `json:"operation" db:"operation"`
	Reason          string    `json:"reason" db:"reason"`
	Notes           string    `json:"notes" db:"notes"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}
