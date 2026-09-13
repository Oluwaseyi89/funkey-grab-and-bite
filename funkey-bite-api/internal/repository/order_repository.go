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

// paymentScanFields holds scan destinations for the nullable payment columns
// shared by every query that reads a full order row, so each one doesn't have
// to repeat six sql.Null* unwrap blocks by hand.
type paymentScanFields struct {
	reference     sql.NullString
	accountNumber sql.NullString
	accountName   sql.NullString
	bankName      sql.NullString
	expiresAt     sql.NullTime
	paidAt        sql.NullTime
}

func (p *paymentScanFields) args() []any {
	return []any{
		&p.reference,
		&p.accountNumber,
		&p.accountName,
		&p.bankName,
		&p.expiresAt,
		&p.paidAt,
	}
}

func (p *paymentScanFields) applyTo(order *models.Order) {
	if p.reference.Valid {
		order.PaymentReference = &p.reference.String
	}
	if p.accountNumber.Valid {
		order.PaymentAccountNumber = &p.accountNumber.String
	}
	if p.accountName.Valid {
		order.PaymentAccountName = &p.accountName.String
	}
	if p.bankName.Valid {
		order.PaymentBankName = &p.bankName.String
	}
	if p.expiresAt.Valid {
		order.PaymentExpiresAt = &p.expiresAt.Time
	}
	if p.paidAt.Valid {
		order.PaymentPaidAt = &p.paidAt.Time
	}
}

func (r *OrderRepository) Create(order *models.Order) (*models.Order, error) {
	query := `
		INSERT INTO orders (
			order_number, user_id, customer_name, customer_phone,
			customer_email, order_type, status, total_amount,
			notes, pickup_time, created_at, payment_method, payment_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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
		order.PaymentMethod,
		order.PaymentStatus,
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
	orderQuery := `
		SELECT id, order_number, user_id, customer_name, customer_phone,
		       customer_email, order_type, status, total_amount, notes,
		       pickup_time, created_at, payment_method, payment_status,
		       payment_reference, payment_account_number, payment_account_name,
		       payment_bank_name, payment_expires_at, payment_paid_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	var userID sql.NullInt64
	var pickupTime sql.NullTime
	var paymentFields paymentScanFields

	dest := append([]any{
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
		&order.PaymentMethod,
		&order.PaymentStatus,
	}, paymentFields.args()...)

	err := r.db.QueryRow(orderQuery, id).Scan(dest...)

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
	paymentFields.applyTo(&order)

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
		       pickup_time, created_at, payment_method, payment_status,
		       payment_reference, payment_account_number, payment_account_name,
		       payment_bank_name, payment_expires_at, payment_paid_at
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
		var paymentFields paymentScanFields

		dest := append([]any{
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
			&order.PaymentMethod,
			&order.PaymentStatus,
		}, paymentFields.args()...)

		err := rows.Scan(dest...)
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
		paymentFields.applyTo(&order)

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

	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read updated order status rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
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

func (r *OrderRepository) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
	query := `
        SELECT id, order_number, user_id, customer_name, customer_phone,
               customer_email, order_type, status, total_amount, notes,
               pickup_time, created_at, payment_method, payment_status,
               payment_reference, payment_account_number, payment_account_name,
               payment_bank_name, payment_expires_at, payment_paid_at
        FROM orders
        WHERE order_number = $1
    `

	var order models.Order
	var userID sql.NullInt64
	var pickupTime sql.NullTime
	var paymentFields paymentScanFields

	dest := append([]any{
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
		&order.PaymentMethod,
		&order.PaymentStatus,
	}, paymentFields.args()...)

	row := r.db.QueryRow(query, orderNumber)
	err := row.Scan(dest...)

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
	paymentFields.applyTo(&order)

	items, err := r.getOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

func (r *OrderRepository) GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error) {
	query := `
        SELECT id, order_number, user_id, customer_name, customer_phone,
               customer_email, order_type, status, total_amount, notes,
               pickup_time, estimated_ready_time, created_at, payment_method,
               payment_status, payment_reference, payment_account_number,
               payment_account_name, payment_bank_name, payment_expires_at,
               payment_paid_at
        FROM orders
        WHERE order_number = $1 AND customer_phone = $2
    `

	var order models.Order
	var userID sql.NullInt64
	var pickupTime sql.NullTime
	var estimatedReadyTime sql.NullTime
	var paymentFields paymentScanFields

	dest := append([]any{
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
		&estimatedReadyTime,
		&order.CreatedAt,
		&order.PaymentMethod,
		&order.PaymentStatus,
	}, paymentFields.args()...)

	row := r.db.QueryRow(query, orderNumber, phone)
	err := row.Scan(dest...)

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
	if estimatedReadyTime.Valid {
		order.EstimatedReadyTime = &estimatedReadyTime.Time
	}
	paymentFields.applyTo(&order)

	items, err := r.getOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

func (r *OrderRepository) CancelOrder(id int) error {
	query := `
        UPDATE orders 
        SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND status = 'pending'
        RETURNING id
    `

	var updatedID int
	err := r.db.QueryRow(query, id).Scan(&updatedID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("order not found or cannot be cancelled (already processed)")
	}
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

func (r *OrderRepository) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *OrderRepository) CreateOrderWithTransaction(tx *sql.Tx, order *models.Order) (*models.Order, error) {
	query := `
        INSERT INTO orders (
            order_number, user_id, customer_name, customer_phone,
            customer_email, order_type, status, total_amount,
            notes, pickup_time, created_at, payment_method, payment_status
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id, created_at
    `

	err := tx.QueryRow(
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
		order.PaymentMethod,
		order.PaymentStatus,
	).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) CreateOrderItemWithTransaction(tx *sql.Tx, item *models.OrderItem) (*models.OrderItem, error) {
	query := `
        INSERT INTO order_items (
            order_id, menu_item_id, name, quantity, unit_price, special_instructions
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	err := tx.QueryRow(
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

func (r *OrderRepository) GetOrderByPaymentReference(reference string) (*models.Order, error) {
	query := `
        SELECT id, order_number, user_id, customer_name, customer_phone,
               customer_email, order_type, status, total_amount, notes,
               pickup_time, created_at, payment_method, payment_status,
               payment_reference, payment_account_number, payment_account_name,
               payment_bank_name, payment_expires_at, payment_paid_at
        FROM orders
        WHERE payment_reference = $1
    `

	var order models.Order
	var userID sql.NullInt64
	var pickupTime sql.NullTime
	var paymentFields paymentScanFields

	dest := append([]any{
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
		&order.PaymentMethod,
		&order.PaymentStatus,
	}, paymentFields.args()...)

	row := r.db.QueryRow(query, reference)
	err := row.Scan(dest...)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order by payment reference: %w", err)
	}

	if userID.Valid {
		val := int(userID.Int64)
		order.UserID = &val
	}
	if pickupTime.Valid {
		order.PickupTime = &pickupTime.Time
	}
	paymentFields.applyTo(&order)

	items, err := r.getOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

func (r *OrderRepository) SetOrderPaymentInitiated(orderID int, reference, accountNumber, accountName, bankName string, expiresAt *time.Time) error {
	query := `
        UPDATE orders
        SET payment_status = 'pending',
            payment_reference = $1,
            payment_account_number = $2,
            payment_account_name = $3,
            payment_bank_name = $4,
            payment_expires_at = $5
        WHERE id = $6
    `

	_, err := r.db.Exec(query, reference, accountNumber, accountName, bankName, expiresAt, orderID)
	if err != nil {
		return fmt.Errorf("failed to record payment initiation: %w", err)
	}

	return nil
}

func (r *OrderRepository) MarkOrderPaymentPaid(orderID int, paidAt time.Time) error {
	query := `
        UPDATE orders
        SET payment_status = 'paid', payment_paid_at = $1
        WHERE id = $2
    `

	_, err := r.db.Exec(query, paidAt, orderID)
	if err != nil {
		return fmt.Errorf("failed to mark order payment paid: %w", err)
	}

	return nil
}
