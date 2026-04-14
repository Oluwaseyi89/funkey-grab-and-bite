package repository

import (
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetByID(id int) (*models.InventoryItem, error) {
	query := `
        SELECT id, menu_item_id, name, current_stock, minimum_stock, reorder_point,
               unit, is_active, last_restocked, created_at, updated_at
        FROM inventory_items
        WHERE id = $1
    `

	var item models.InventoryItem
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&item.ID,
		&item.MenuItemID,
		&item.Name,
		&item.CurrentStock,
		&item.MinimumStock,
		&item.ReorderPoint,
		&item.Unit,
		&item.IsActive,
		&item.LastRestocked,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return &item, nil
}

func (r *InventoryRepository) GetByMenuItemID(menuItemID int) (*models.InventoryItem, error) {
	query := `
        SELECT id, menu_item_id, name, current_stock, minimum_stock, reorder_point,
               unit, is_active, last_restocked, created_at, updated_at
        FROM inventory_items
        WHERE menu_item_id = $1
    `

	var item models.InventoryItem
	row := r.db.QueryRow(query, menuItemID)

	err := row.Scan(
		&item.ID,
		&item.MenuItemID,
		&item.Name,
		&item.CurrentStock,
		&item.MinimumStock,
		&item.ReorderPoint,
		&item.Unit,
		&item.IsActive,
		&item.LastRestocked,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return &item, nil
}

func (r *InventoryRepository) GetAll() ([]models.InventoryItem, error) {
	query := `
        SELECT id, menu_item_id, name, current_stock, minimum_stock, reorder_point,
               unit, is_active, last_restocked, created_at, updated_at
        FROM inventory_items
        ORDER BY name
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory items: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.MenuItemID,
			&item.Name,
			&item.CurrentStock,
			&item.MinimumStock,
			&item.ReorderPoint,
			&item.Unit,
			&item.IsActive,
			&item.LastRestocked,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *InventoryRepository) GetLowStock() ([]models.InventoryItem, error) {
	query := `
        SELECT id, menu_item_id, name, current_stock, minimum_stock, reorder_point,
               unit, is_active, last_restocked, created_at, updated_at
        FROM inventory_items
        WHERE current_stock <= reorder_point AND is_active = true
        ORDER BY current_stock ASC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.MenuItemID,
			&item.Name,
			&item.CurrentStock,
			&item.MinimumStock,
			&item.ReorderPoint,
			&item.Unit,
			&item.IsActive,
			&item.LastRestocked,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *InventoryRepository) UpdateStock(itemID int, newStock int, operation string, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var currentStock int
	err = tx.QueryRow("SELECT current_stock FROM inventory_items WHERE id = $1", itemID).Scan(&currentStock)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	query := `
        UPDATE inventory_items 
        SET current_stock = $1, 
            last_restocked = CASE WHEN $2 = 'add' THEN $3 ELSE last_restocked END,
            updated_at = $3
        WHERE id = $4
    `

	now := time.Now()
	_, err = tx.Exec(query, newStock, operation, now, itemID)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	historyQuery := `
        INSERT INTO inventory_history 
        (inventory_item_id, previous_stock, new_stock, change, operation, reason, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	change := newStock - currentStock
	_, err = tx.Exec(historyQuery, itemID, currentStock, newStock, change, operation, reason, now)
	if err != nil {
		return fmt.Errorf("failed to record inventory history: %w", err)
	}

	if newStock <= 0 {
		alertQuery := `
            INSERT INTO inventory_alerts (inventory_item_id, alert_type, message, created_at)
            VALUES ($1, 'out_of_stock', $2, $3)
            ON CONFLICT (inventory_item_id, alert_type) 
            DO UPDATE SET is_resolved = false, created_at = $3
            WHERE inventory_alerts.is_resolved = true
        `
		message := fmt.Sprintf("Item %d is out of stock", itemID)
		tx.Exec(alertQuery, itemID, message, now)
	} else if newStock <= 10 {
		alertQuery := `
            INSERT INTO inventory_alerts (inventory_item_id, alert_type, message, created_at)
            VALUES ($1, 'low_stock', $2, $3)
            ON CONFLICT (inventory_item_id, alert_type) 
            DO UPDATE SET is_resolved = false, created_at = $3
            WHERE inventory_alerts.is_resolved = true
        `
		message := fmt.Sprintf("Item %d is running low: %d remaining", itemID, newStock)
		tx.Exec(alertQuery, itemID, message, now)
	}

	return tx.Commit()
}

func (r *InventoryRepository) AdjustStockByMenuItemID(menuItemID int, quantity int, operation string, reason string) (*models.InventoryItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var itemID int
	var currentStock int
	err = tx.QueryRow(
		`SELECT id, current_stock FROM inventory_items WHERE menu_item_id = $1 FOR UPDATE`,
		menuItemID,
	).Scan(&itemID, &currentStock)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to lock inventory item: %w", err)
	}

	var newStock int
	switch operation {
	case "add":
		newStock = currentStock + quantity
	case "subtract":
		newStock = currentStock - quantity
		if newStock < 0 {
			return nil, fmt.Errorf("insufficient stock. Current: %d, Requested: %d", currentStock, quantity)
		}
	case "set":
		newStock = quantity
		if newStock < 0 {
			return nil, fmt.Errorf("stock cannot be negative")
		}
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}

	now := time.Now()
	_, err = tx.Exec(
		`UPDATE inventory_items
		 SET current_stock = $1,
		     last_restocked = CASE WHEN $2 = 'add' THEN $3 ELSE last_restocked END,
		     updated_at = $3
		 WHERE id = $4`,
		newStock, operation, now, itemID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO inventory_history
		 (inventory_item_id, previous_stock, new_stock, change, operation, reason, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		itemID, currentStock, newStock, newStock-currentStock, operation, reason, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to record inventory history: %w", err)
	}

	if newStock <= 0 {
		message := fmt.Sprintf("Item %d is out of stock", itemID)
		_, _ = tx.Exec(
			`INSERT INTO inventory_alerts (inventory_item_id, alert_type, message, created_at)
			 VALUES ($1, 'out_of_stock', $2, $3)
			 ON CONFLICT (inventory_item_id, alert_type)
			 DO UPDATE SET is_resolved = false, created_at = $3
			 WHERE inventory_alerts.is_resolved = true`,
			itemID, message, now,
		)
	} else if newStock <= 10 {
		message := fmt.Sprintf("Item %d is running low: %d remaining", itemID, newStock)
		_, _ = tx.Exec(
			`INSERT INTO inventory_alerts (inventory_item_id, alert_type, message, created_at)
			 VALUES ($1, 'low_stock', $2, $3)
			 ON CONFLICT (inventory_item_id, alert_type)
			 DO UPDATE SET is_resolved = false, created_at = $3
			 WHERE inventory_alerts.is_resolved = true`,
			itemID, message, now,
		)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit inventory update: %w", err)
	}
	committed = true

	return r.GetByID(itemID)
}

func (r *InventoryRepository) CreateInventoryItem(item *models.InventoryItem) (*models.InventoryItem, error) {
	query := `
        INSERT INTO inventory_items 
        (menu_item_id, name, current_stock, minimum_stock, reorder_point, unit, is_active, last_restocked, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `

	now := time.Now()
	err := r.db.QueryRow(
		query,
		item.MenuItemID,
		item.Name,
		item.CurrentStock,
		item.MinimumStock,
		item.ReorderPoint,
		item.Unit,
		item.IsActive,
		item.LastRestocked,
		now,
		now,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	return item, nil
}

func (r *InventoryRepository) GetAlerts(resolved bool) ([]models.InventoryAlert, error) {
	query := `
        SELECT id, inventory_item_id, alert_type, message, is_resolved, created_at, resolved_at
        FROM inventory_alerts
        WHERE is_resolved = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}
	defer rows.Close()

	var alerts []models.InventoryAlert
	for rows.Next() {
		var alert models.InventoryAlert
		var resolvedAt sql.NullTime
		err := rows.Scan(
			&alert.ID,
			&alert.InventoryItemID,
			&alert.AlertType,
			&alert.Message,
			&alert.IsResolved,
			&alert.CreatedAt,
			&resolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}
		if resolvedAt.Valid {
			alert.ResolvedAt = &resolvedAt.Time
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *InventoryRepository) ResolveAlert(alertID int) error {
	query := `UPDATE inventory_alerts SET is_resolved = true, resolved_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), alertID)
	return err
}
