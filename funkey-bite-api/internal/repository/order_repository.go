package repository

import (
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) (*models.Order, error) {
	query := `
		INSERT INTO orders (
			order_number, user_id, customer_name, customer_phone, 
			customer_email, order_type, status, total_amount, 
			notes, pickup_time, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		order.OrderNumber,
		order.UserID,
		order.CustomerName,
		order.CustomerPhone,
		order.CustomerEmail,
		order.OrderType,
		order.Status,
		order.TotalAmount,
		order.Notes,
		order.PickupTime,
		order.CreatedAt,
	).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) (*models.OrderItem, error) {
	query := `
		INSERT INTO order_items (
			order_id, menu_item_id, name, quantity, unit_price, special_instructions
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		item.OrderID,
		item.MenuItemID,
		item.Name,
		item.Quantity,
		item.UnitPrice,
		item.SpecialInstructions,
	).Scan(&item.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create order item: %w", err)
	}

	return item, nil
}

func (r *OrderRepository) GetOrderWithItems(id int) (*models.Order, error) {
	// First, get the order
	orderQuery := `
		SELECT id, order_number, user_id, customer_name, customer_phone, 
		       customer_email, order_type, status, total_amount, notes, 
		       pickup_time, created_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	var userID sql.NullInt64
	var pickupTime sql.NullTime

	err := r.db.QueryRow(orderQuery, id).Scan(
		&order.ID,
		&order.OrderNumber,
		&userID,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.CustomerEmail,
		&order.OrderType,
		&order.Status,
		&order.TotalAmount,
		&order.Notes,
		&pickupTime,
		&order.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if userID.Valid {
		val := int(userID.Int64)
		order.UserID = &val
	}
	if pickupTime.Valid {
		order.PickupTime = &pickupTime.Time
	}

	// Get order items
	itemsQuery := `
		SELECT id, order_id, menu_item_id, name, quantity, unit_price, special_instructions
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var specialInstructions sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Name,
			&item.Quantity,
			&item.UnitPrice,
			&specialInstructions,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if specialInstructions.Valid {
			item.SpecialInstructions = &specialInstructions.String
		}

		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID int) ([]models.Order, error) {
	query := `
		SELECT id, order_number, user_id, customer_name, customer_phone, 
		       customer_email, order_type, status, total_amount, notes, 
		       pickup_time, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var uid sql.NullInt64
		var pickupTime sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&uid,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CustomerEmail,
			&order.OrderType,
			&order.Status,
			&order.TotalAmount,
			&order.Notes,
			&pickupTime,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if uid.Valid {
			val := int(uid.Int64)
			order.UserID = &val
		}
		if pickupTime.Valid {
			order.PickupTime = &pickupTime.Time
		}

		// Get items for this order
		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(id int, status string) error {
	query := `
		UPDATE orders 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (r *OrderRepository) getOrderItems(orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, menu_item_id, name, quantity, unit_price, special_instructions
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var specialInstructions sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Name,
			&item.Quantity,
			&item.UnitPrice,
			&specialInstructions,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if specialInstructions.Valid {
			item.SpecialInstructions = &specialInstructions.String
		}

		items = append(items, item)
	}

	return items, nil
}
