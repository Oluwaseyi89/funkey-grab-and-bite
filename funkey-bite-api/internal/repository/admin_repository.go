package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// Dashboard Statistics
// func (r *AdminRepository) GetDashboardStats(fromDate, toDate time.Time) (*models.AdminStats, error) {
// 	var stats models.AdminStats

// 	// Total orders
// 	query := `SELECT COUNT(*) FROM orders WHERE created_at BETWEEN $1 AND $2`
// 	err := r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalOrders)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get total orders: %w", err)
// 	}

// 	// Total revenue
// 	query = `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at BETWEEN $1 AND $2`
// 	err = r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalRevenue)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get total revenue: %w", err)
// 	}

// 	// Pending orders
// 	query = `SELECT COUNT(*) FROM orders WHERE status = 'pending'`
// 	err = r.db.QueryRow(query).Scan(&stats.PendingOrders)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get pending orders: %w", err)
// 	}

// 	// Active catering requests
// 	query = `SELECT COUNT(*) FROM catering_requests WHERE status IN ('pending', 'confirmed')`
// 	err = r.db.QueryRow(query).Scan(&stats.ActiveCatering)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get active catering: %w", err)
// 	}

// 	// New customers (last 7 days)
// 	query = `SELECT COUNT(*) FROM users WHERE created_at >= $1`
// 	err = r.db.QueryRow(query, time.Now().AddDate(0, 0, -7)).Scan(&stats.NewCustomers)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get new customers: %w", err)
// 	}

// 	// Popular items
// 	query := `
// 		SELECT oi.menu_item_id, mi.name,
// 		       SUM(oi.quantity) as total_sold,
// 		       SUM(oi.quantity * oi.unit_price) as revenue
// 		FROM order_items oi
// 		JOIN menu_items mi ON oi.menu_item_id = mi.id
// 		JOIN orders o ON oi.order_id = o.id
// 		WHERE o.created_at BETWEEN $1 AND $2
// 		GROUP BY oi.menu_item_id, mi.name
// 		ORDER BY total_sold DESC
// 		LIMIT 5
// 	`
// 	rows, err := r.db.Query(query, fromDate, toDate)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get popular items: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var item models.MenuItemStats
// 		err := rows.Scan(&item.ID, &item.Name, &item.TotalSold, &item.Revenue)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan popular item: %w", err)
// 		}
// 		stats.PopularItems = append(stats.PopularItems, item)
// 	}

// 	return &stats, nil
// }

func (r *AdminRepository) GetDashboardStats(fromDate, toDate time.Time) (*models.AdminStats, error) {
	var stats models.AdminStats
	var query string

	// Total orders
	query = `SELECT COUNT(*) FROM orders WHERE created_at BETWEEN $1 AND $2`
	err := r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}

	// Total revenue
	query = `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at BETWEEN $1 AND $2`
	err = r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}

	// Pending orders
	query = `SELECT COUNT(*) FROM orders WHERE status = 'pending'`
	err = r.db.QueryRow(query).Scan(&stats.PendingOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}

	// Active catering requests
	query = `SELECT COUNT(*) FROM catering_requests WHERE status IN ('pending', 'confirmed')`
	err = r.db.QueryRow(query).Scan(&stats.ActiveCatering)
	if err != nil {
		return nil, fmt.Errorf("failed to get active catering: %w", err)
	}

	// New customers (last 7 days)
	query = `SELECT COUNT(*) FROM users WHERE created_at >= $1`
	err = r.db.QueryRow(query, time.Now().AddDate(0, 0, -7)).Scan(&stats.NewCustomers)
	if err != nil {
		return nil, fmt.Errorf("failed to get new customers: %w", err)
	}

	// Popular items
	query = `
		SELECT oi.menu_item_id, mi.name, 
		       SUM(oi.quantity) as total_sold,
		       SUM(oi.quantity * oi.unit_price) as revenue
		FROM order_items oi
		JOIN menu_items mi ON oi.menu_item_id = mi.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY oi.menu_item_id, mi.name
		ORDER BY total_sold DESC
		LIMIT 5
	`
	rows, err := r.db.Query(query, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.MenuItemStats
		err := rows.Scan(&item.ID, &item.Name, &item.TotalSold, &item.Revenue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan popular item: %w", err)
		}
		stats.PopularItems = append(stats.PopularItems, item)
	}

	return &stats, nil
}

// Sales Reports
func (r *AdminRepository) GetSalesReport(fromDate, toDate time.Time) ([]models.SalesReport, error) {
	query := `
		SELECT DATE(created_at) as date,
		       COUNT(*) as total_orders,
		       COALESCE(SUM(total_amount), 0) as total_revenue,
		       COALESCE(AVG(total_amount), 0) as avg_order
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	rows, err := r.db.Query(query, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales report: %w", err)
	}
	defer rows.Close()

	var reports []models.SalesReport
	for rows.Next() {
		var report models.SalesReport
		err := rows.Scan(&report.Date, &report.TotalOrders, &report.TotalRevenue, &report.AverageOrder)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales report: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// Menu Management
func (r *AdminRepository) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	tagsJSON, _ := json.Marshal(item.Tags)
	var nutritionalInfoJSON []byte
	if item.NutritionalInfo != nil {
		nutritionalInfoJSON, _ = json.Marshal(item.NutritionalInfo)
	}

	query := `
		INSERT INTO menu_items (
			category_id, name, description, price, image_url,
			is_available, is_pre_order, preparation_time, tags, nutritional_info
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	var createdAt time.Time
	err := r.db.QueryRow(
		query,
		item.CategoryID,
		item.Name,
		item.Description,
		item.Price,
		item.ImageURL,
		item.IsAvailable,
		item.IsPreOrder,
		item.PreparationTime,
		tagsJSON,
		nutritionalInfoJSON,
	).Scan(&item.ID, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	item.CreatedAt = createdAt
	return item, nil
}

func (r *AdminRepository) UpdateMenuItem(item *models.MenuItem) error {
	tagsJSON, _ := json.Marshal(item.Tags)
	var nutritionalInfoJSON []byte
	if item.NutritionalInfo != nil {
		nutritionalInfoJSON, _ = json.Marshal(item.NutritionalInfo)
	}

	query := `
		UPDATE menu_items SET
			category_id = $1,
			name = $2,
			description = $3,
			price = $4,
			image_url = $5,
			is_available = $6,
			is_pre_order = $7,
			preparation_time = $8,
			tags = $9,
			nutritional_info = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`

	_, err := r.db.Exec(
		query,
		item.CategoryID,
		item.Name,
		item.Description,
		item.Price,
		item.ImageURL,
		item.IsAvailable,
		item.IsPreOrder,
		item.PreparationTime,
		tagsJSON,
		nutritionalInfoJSON,
		item.ID,
	)

	return err
}

func (r *AdminRepository) DeleteMenuItem(id int) error {
	query := `DELETE FROM menu_items WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Order Management
func (r *AdminRepository) GetAllOrders(limit, offset int, status string) ([]models.Order, error) {
	query := `
		SELECT id, order_number, user_id, customer_name, customer_phone, 
		       customer_email, order_type, status, total_amount, notes, 
		       pickup_time, created_at
		FROM orders
		WHERE ($1 = '' OR status = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var userID sql.NullInt64
		var pickupTime sql.NullTime

		err := rows.Scan(
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
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if userID.Valid {
			val := int(userID.Int64)
			order.UserID = &val
		}
		if pickupTime.Valid {
			order.PickupTime = &pickupTime.Time
		}

		// Get order items
		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *AdminRepository) getOrderItems(orderID int) ([]models.OrderItem, error) {
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

func (r *AdminRepository) GetOrdersCount(status string) (int, error) {
	query := `SELECT COUNT(*) FROM orders WHERE ($1 = '' OR status = $1)`
	var count int
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get orders count: %w", err)
	}
	return count, nil
}

// User Management
func (r *AdminRepository) GetAllUsers(limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, phone, email, full_name, is_verified, is_active, 
		       last_login, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var email sql.NullString
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.Phone,
			&email,
			&user.FullName,
			&user.IsVerified,
			&user.IsActive,
			&lastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if email.Valid {
			user.Email = &email.String
		}
		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *AdminRepository) GetUsersCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get users count: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) UpdateUserStatus(userID int, isActive bool) error {
	query := `UPDATE users SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, isActive, userID)
	return err
}

// package repository

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
// )

// type AdminRepository struct {
// 	db *sql.DB
// }

// func NewAdminRepository(db *sql.DB) *AdminRepository {
// 	return &AdminRepository{db: db}
// }

// // Dashboard Statistics
// func (r *AdminRepository) GetDashboardStats(fromDate, toDate time.Time) (*models.AdminStats, error) {
// 	var stats models.AdminStats

// 	// Total orders
// 	query := `SELECT COUNT(*) FROM orders WHERE created_at BETWEEN $1 AND $2`
// 	err := r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalOrders)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get total orders: %w", err)
// 	}

// 	// Total revenue
// 	query = `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at BETWEEN $1 AND $2`
// 	err = r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalRevenue)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get total revenue: %w", err)
// 	}

// 	// Pending orders
// 	query = `SELECT COUNT(*) FROM orders WHERE status = 'pending'`
// 	err = r.db.QueryRow(query).Scan(&stats.PendingOrders)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get pending orders: %w", err)
// 	}

// 	// Active catering requests
// 	query = `SELECT COUNT(*) FROM catering_requests WHERE status IN ('pending', 'confirmed')`
// 	err = r.db.QueryRow(query).Scan(&stats.ActiveCatering)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get active catering: %w", err)
// 	}

// 	// New customers (last 7 days)
// 	query = `SELECT COUNT(*) FROM users WHERE created_at >= $1`
// 	err = r.db.QueryRow(query, time.Now().AddDate(0, 0, -7)).Scan(&stats.NewCustomers)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get new customers: %w", err)
// 	}

// 	// Popular items
// 	query = `
// 		SELECT oi.menu_item_id, mi.name,
// 		       SUM(oi.quantity) as total_sold,
// 		       SUM(oi.quantity * oi.unit_price) as revenue
// 		FROM order_items oi
// 		JOIN menu_items mi ON oi.menu_item_id = mi.id
// 		JOIN orders o ON oi.order_id = o.id
// 		WHERE o.created_at BETWEEN $1 AND $2
// 		GROUP BY oi.menu_item_id, mi.name
// 		ORDER BY total_sold DESC
// 		LIMIT 5
// 	`
// 	rows, err := r.db.Query(query, fromDate, toDate)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get popular items: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var item models.MenuItemStats
// 		err := rows.Scan(&item.ID, &item.Name, &item.TotalSold, &item.Revenue)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan popular item: %w", err)
// 		}
// 		stats.PopularItems = append(stats.PopularItems, item)
// 	}

// 	return &stats, nil
// }

// // Sales Reports
// func (r *AdminRepository) GetSalesReport(fromDate, toDate time.Time) ([]models.SalesReport, error) {
// 	query := `
// 		SELECT DATE(created_at) as date,
// 		       COUNT(*) as total_orders,
// 		       COALESCE(SUM(total_amount), 0) as total_revenue,
// 		       COALESCE(AVG(total_amount), 0) as avg_order
// 		FROM orders
// 		WHERE created_at BETWEEN $1 AND $2
// 		GROUP BY DATE(created_at)
// 		ORDER BY date
// 	`

// 	rows, err := r.db.Query(query, fromDate, toDate)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get sales report: %w", err)
// 	}
// 	defer rows.Close()

// 	var reports []models.SalesReport
// 	for rows.Next() {
// 		var report models.SalesReport
// 		err := rows.Scan(&report.Date, &report.TotalOrders, &report.TotalRevenue, &report.AverageOrder)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan sales report: %w", err)
// 		}
// 		reports = append(reports, report)
// 	}

// 	return reports, nil
// }

// // Menu Management
// func (r *AdminRepository) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
// 	tagsJSON, _ := json.Marshal(item.Tags)
// 	var nutritionalInfoJSON []byte
// 	if item.NutritionalInfo != nil {
// 		nutritionalInfoJSON, _ = json.Marshal(item.NutritionalInfo)
// 	}

// 	query := `
// 		INSERT INTO menu_items (
// 			category_id, name, description, price, image_url,
// 			is_available, is_pre_order, preparation_time, tags, nutritional_info
// 		)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
// 		RETURNING id, created_at
// 	`

// 	var createdAt time.Time
// 	err := r.db.QueryRow(
// 		query,
// 		item.CategoryID,
// 		item.Name,
// 		item.Description,
// 		item.Price,
// 		item.ImageURL,
// 		item.IsAvailable,
// 		item.IsPreOrder,
// 		item.PreparationTime,
// 		tagsJSON,
// 		nutritionalInfoJSON,
// 	).Scan(&item.ID, &createdAt)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create menu item: %w", err)
// 	}

// 	item.CreatedAt = createdAt
// 	return item, nil
// }

// func (r *AdminRepository) UpdateMenuItem(item *models.MenuItem) error {
// 	tagsJSON, _ := json.Marshal(item.Tags)
// 	var nutritionalInfoJSON []byte
// 	if item.NutritionalInfo != nil {
// 		nutritionalInfoJSON, _ = json.Marshal(item.NutritionalInfo)
// 	}

// 	query := `
// 		UPDATE menu_items SET
// 			category_id = $1,
// 			name = $2,
// 			description = $3,
// 			price = $4,
// 			image_url = $5,
// 			is_available = $6,
// 			is_pre_order = $7,
// 			preparation_time = $8,
// 			tags = $9,
// 			nutritional_info = $10,
// 			updated_at = CURRENT_TIMESTAMP
// 		WHERE id = $11
// 	`

// 	_, err := r.db.Exec(
// 		query,
// 		item.CategoryID,
// 		item.Name,
// 		item.Description,
// 		item.Price,
// 		item.ImageURL,
// 		item.IsAvailable,
// 		item.IsPreOrder,
// 		item.PreparationTime,
// 		tagsJSON,
// 		nutritionalInfoJSON,
// 		item.ID,
// 	)

// 	return err
// }

// func (r *AdminRepository) DeleteMenuItem(id int) error {
// 	query := `DELETE FROM menu_items WHERE id = $1`
// 	_, err := r.db.Exec(query, id)
// 	return err
// }

// // Order Management
// func (r *AdminRepository) GetAllOrders(limit, offset int, status string) ([]models.Order, error) {
// 	query := `
// 		SELECT id, order_number, user_id, customer_name, customer_phone,
// 		       customer_email, order_type, status, total_amount, notes,
// 		       pickup_time, created_at
// 		FROM orders
// 		WHERE ($1 = '' OR status = $1)
// 		ORDER BY created_at DESC
// 		LIMIT $2 OFFSET $3
// 	`

// 	rows, err := r.db.Query(query, status, limit, offset)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get orders: %w", err)
// 	}
// 	defer rows.Close()

// 	var orders []models.Order
// 	for rows.Next() {
// 		var order models.Order
// 		var userID sql.NullInt64
// 		var pickupTime sql.NullTime

// 		err := rows.Scan(
// 			&order.ID,
// 			&order.OrderNumber,
// 			&userID,
// 			&order.CustomerName,
// 			&order.CustomerPhone,
// 			&order.CustomerEmail,
// 			&order.OrderType,
// 			&order.Status,
// 			&order.TotalAmount,
// 			&order.Notes,
// 			&pickupTime,
// 			&order.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan order: %w", err)
// 		}

// 		if userID.Valid {
// 			val := int(userID.Int64)
// 			order.UserID = &val
// 		}
// 		if pickupTime.Valid {
// 			order.PickupTime = &pickupTime.Time
// 		}

// 		// Get order items
// 		items, err := r.getOrderItems(order.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		order.Items = items

// 		orders = append(orders, order)
// 	}

// 	return orders, nil
// }

// func (r *AdminRepository) getOrderItems(orderID int) ([]models.OrderItem, error) {
// 	query := `
// 		SELECT id, order_id, menu_item_id, name, quantity, unit_price, special_instructions
// 		FROM order_items
// 		WHERE order_id = $1
// 		ORDER BY id
// 	`

// 	rows, err := r.db.Query(query, orderID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get order items: %w", err)
// 	}
// 	defer rows.Close()

// 	var items []models.OrderItem
// 	for rows.Next() {
// 		var item models.OrderItem
// 		var specialInstructions sql.NullString

// 		err := rows.Scan(
// 			&item.ID,
// 			&item.OrderID,
// 			&item.MenuItemID,
// 			&item.Name,
// 			&item.Quantity,
// 			&item.UnitPrice,
// 			&specialInstructions,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan order item: %w", err)
// 		}

// 		if specialInstructions.Valid {
// 			item.SpecialInstructions = &specialInstructions.String
// 		}

// 		items = append(items, item)
// 	}

// 	return items, nil
// }

// // User Management
// func (r *AdminRepository) GetAllUsers(limit, offset int) ([]models.User, error) {
// 	query := `
// 		SELECT id, phone, email, full_name, is_verified, is_active,
// 		       last_login, created_at, updated_at
// 		FROM users
// 		ORDER BY created_at DESC
// 		LIMIT $1 OFFSET $2
// 	`

// 	rows, err := r.db.Query(query, limit, offset)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get users: %w", err)
// 	}
// 	defer rows.Close()

// 	var users []models.User
// 	for rows.Next() {
// 		var user models.User
// 		var email sql.NullString
// 		var lastLogin sql.NullTime

// 		err := rows.Scan(
// 			&user.ID,
// 			&user.Phone,
// 			&email,
// 			&user.FullName,
// 			&user.IsVerified,
// 			&user.IsActive,
// 			&lastLogin,
// 			&user.CreatedAt,
// 			&user.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan user: %w", err)
// 		}

// 		if email.Valid {
// 			user.Email = &email.String
// 		}
// 		if lastLogin.Valid {
// 			user.LastLogin = &lastLogin.Time
// 		}

// 		users = append(users, user)
// 	}

// 	return users, nil
// }

// func (r *AdminRepository) UpdateUserStatus(userID int, isActive bool) error {
// 	query := `UPDATE users SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
// 	_, err := r.db.Exec(query, isActive, userID)
// 	return err
// }
