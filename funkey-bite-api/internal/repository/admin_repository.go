package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetDashboardStats(fromDate, toDate time.Time) (*models.AdminStats, error) {
	var stats models.AdminStats
	var query string

	query = `SELECT COUNT(*) FROM orders WHERE created_at BETWEEN $1 AND $2`
	err := r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}

	query = `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE created_at BETWEEN $1 AND $2`
	err = r.db.QueryRow(query, fromDate, toDate).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}

	query = `SELECT COUNT(*) FROM orders WHERE status = 'pending'`
	err = r.db.QueryRow(query).Scan(&stats.PendingOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders: %w", err)
	}

	query = `SELECT COUNT(*) FROM catering_requests WHERE status IN ('pending', 'confirmed')`
	err = r.db.QueryRow(query).Scan(&stats.ActiveCatering)
	if err != nil {
		return nil, fmt.Errorf("failed to get active catering: %w", err)
	}

	query = `SELECT COUNT(*) FROM users WHERE created_at >= $1`
	err = r.db.QueryRow(query, time.Now().AddDate(0, 0, -7)).Scan(&stats.NewCustomers)
	if err != nil {
		return nil, fmt.Errorf("failed to get new customers: %w", err)
	}

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

func (r *AdminRepository) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	tagsJSON, _ := json.Marshal(item.Tags)
	var nutritionalInfoJSON []byte
	if item.NutritionalInfo != nil {
		nutritionalInfoJSON, _ = json.Marshal(item.NutritionalInfo)
	}

	// Every orderable item needs a stock record, so the item is created with
	// one at zero stock in the same statement; the admin restocks it in
	// Inventory before customers can order it.
	query := `
		WITH new_item AS (
			INSERT INTO menu_items (
				category_id, name, description, price, image_url,
				is_available, is_pre_order, preparation_time, tags, nutritional_info
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, name, created_at
		), new_stock AS (
			INSERT INTO inventory_items (menu_item_id, name, current_stock)
			SELECT id, name, 0 FROM new_item
		)
		SELECT id, created_at FROM new_item
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

func (r *AdminRepository) GetAllOrders(limit, offset int, filter models.OrderListFilter) ([]models.Order, error) {
	query := `
		SELECT id, order_number, user_id, customer_name, customer_phone,
		       customer_email, order_type, status, total_amount, notes,
		       pickup_time, created_at, payment_method, payment_status,
		       payment_reference, payment_account_number, payment_account_name,
		       payment_bank_name, payment_expires_at, payment_paid_at,
		       delivery_address, invoice_number, subtotal_amount, delivery_fee,
		       tax_amount, discount_amount, payment_collected_via
		FROM orders
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR payment_status = $2)
		  AND ($3 = '' OR payment_method = $3)
		  AND ($4::int IS NULL OR user_id = $4)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.db.Query(query, filter.Status, filter.PaymentStatus, filter.PaymentMethod, filter.UserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
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

		err := rows.Scan(dest...)
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

func (r *AdminRepository) GetOrdersCount(filter models.OrderListFilter) (int, error) {
	query := `
		SELECT COUNT(*) FROM orders
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR payment_status = $2)
		  AND ($3 = '' OR payment_method = $3)
		  AND ($4::int IS NULL OR user_id = $4)
	`
	var count int
	err := r.db.QueryRow(query, filter.Status, filter.PaymentStatus, filter.PaymentMethod, filter.UserID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get orders count: %w", err)
	}
	return count, nil
}

// customerSelect lists customers with their order totals. Cancelled orders are
// left out and only paid orders count towards what a customer has spent.
const customerSelect = `
	WITH order_stats AS (
		SELECT user_id,
		       COUNT(*) AS order_count,
		       COALESCE(SUM(total_amount) FILTER (WHERE payment_status = 'paid'), 0) AS total_spent,
		       MIN(created_at) AS first_order_at,
		       MAX(created_at) AS last_order_at
		FROM orders
		WHERE user_id IS NOT NULL AND status <> 'cancelled'
		GROUP BY user_id
	)
	SELECT u.id, u.phone, u.email, u.full_name, u.is_verified, u.is_active,
	       u.last_login, u.created_at, u.updated_at,
	       COALESCE(s.order_count, 0), COALESCE(s.total_spent, 0),
	       s.first_order_at, s.last_order_at
	FROM users u
	LEFT JOIN order_stats s ON s.user_id = u.id
`

// customerWhere applies a CustomerListFilter given as $1 (search) and $2 (status).
const customerWhere = `
	WHERE ($1 = '' OR u.full_name ILIKE $1 OR u.phone ILIKE $1 OR u.email ILIKE $1)
	  AND ($2 = ''
	       OR ($2 = 'active' AND u.is_active)
	       OR ($2 = 'inactive' AND NOT u.is_active)
	       OR ($2 = 'verified' AND u.is_verified)
	       OR ($2 = 'unverified' AND NOT u.is_verified))
`

// customerSearchPattern turns a search into an ILIKE pattern that matches it
// anywhere, treating % and _ as plain characters.
func customerSearchPattern(search string) string {
	if search == "" {
		return ""
	}
	escaped := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(search)
	return "%" + escaped + "%"
}

func scanCustomer(row interface{ Scan(...any) error }) (*models.AdminCustomer, error) {
	var customer models.AdminCustomer
	var email sql.NullString
	var lastLogin, firstOrderAt, lastOrderAt sql.NullTime

	err := row.Scan(
		&customer.ID,
		&customer.Phone,
		&email,
		&customer.FullName,
		&customer.IsVerified,
		&customer.IsActive,
		&lastLogin,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.OrderCount,
		&customer.TotalSpent,
		&firstOrderAt,
		&lastOrderAt,
	)
	if err != nil {
		return nil, err
	}

	if email.Valid {
		customer.Email = &email.String
	}
	if lastLogin.Valid {
		customer.LastLogin = &lastLogin.Time
	}
	if firstOrderAt.Valid {
		customer.FirstOrderAt = &firstOrderAt.Time
	}
	if lastOrderAt.Valid {
		customer.LastOrderAt = &lastOrderAt.Time
	}
	return &customer, nil
}

func (r *AdminRepository) GetAllUsers(limit, offset int, filter models.CustomerListFilter) ([]models.AdminCustomer, error) {
	query := customerSelect + customerWhere + `
		ORDER BY u.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, customerSearchPattern(filter.Search), filter.Status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var customers []models.AdminCustomer
	for rows.Next() {
		customer, err := scanCustomer(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		customers = append(customers, *customer)
	}

	return customers, rows.Err()
}

func (r *AdminRepository) GetUsersCount(filter models.CustomerListFilter) (int, error) {
	query := `SELECT COUNT(*) FROM users u` + customerWhere
	var count int
	err := r.db.QueryRow(query, customerSearchPattern(filter.Search), filter.Status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get users count: %w", err)
	}
	return count, nil
}

// GetCustomer returns one customer with their order totals and the menu
// category they order most. It returns sql.ErrNoRows for an unknown ID.
func (r *AdminRepository) GetCustomer(userID int) (*models.AdminCustomer, error) {
	customer, err := scanCustomer(r.db.QueryRow(customerSelect+` WHERE u.id = $1`, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var category string
	err = r.db.QueryRow(`
		SELECT c.name
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN menu_items m ON m.id = oi.menu_item_id
		JOIN menu_categories c ON c.id = m.category_id
		WHERE o.user_id = $1 AND o.status <> 'cancelled'
		GROUP BY c.name
		ORDER BY SUM(oi.quantity) DESC, c.name
		LIMIT 1
	`, userID).Scan(&category)
	switch {
	case err == nil:
		customer.FavoriteCategory = &category
	case !errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("failed to get favorite category: %w", err)
	}

	return customer, nil
}

func (r *AdminRepository) UpdateUserStatus(userID int, isActive bool) error {
	query := `UPDATE users SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	result, err := r.db.Exec(query, isActive, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AdminRepository) GetAdminByEmail(email string) (*models.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_active, last_login, created_at
		FROM admin_users
		WHERE email = $1
	`

	var admin models.AdminUser
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.PasswordHash,
		&admin.Role,
		&admin.IsActive,
		&lastLogin,
		&admin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin by email: %w", err)
	}

	if lastLogin.Valid {
		admin.LastLogin = &lastLogin.Time
	}

	return &admin, nil
}

func (r *AdminRepository) UpdateAdminLastLogin(adminID int) error {
	query := `UPDATE admin_users SET last_login = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, adminID)
	return err
}

func (r *AdminRepository) GetAdminUsers(limit, offset int) ([]models.AdminUser, error) {
	query := `
		SELECT id, username, email, role, is_active, last_login, created_at
		FROM admin_users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}
	defer rows.Close()

	var admins []models.AdminUser
	for rows.Next() {
		var admin models.AdminUser
		var lastLogin sql.NullTime

		err := rows.Scan(
			&admin.ID,
			&admin.Username,
			&admin.Email,
			&admin.Role,
			&admin.IsActive,
			&lastLogin,
			&admin.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan admin user: %w", err)
		}

		if lastLogin.Valid {
			admin.LastLogin = &lastLogin.Time
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *AdminRepository) CreateAdminUser(admin *models.AdminUser) error {
	query := `
		INSERT INTO admin_users (username, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		admin.Username,
		admin.Email,
		admin.PasswordHash,
		admin.Role,
		admin.IsActive,
	).Scan(&admin.ID, &admin.CreatedAt)

	return err
}

func (r *AdminRepository) UpdateAdminUser(admin *models.AdminUser) error {
	query := `
		UPDATE admin_users SET
			username = $1,
			email = $2,
			role = $3,
			is_active = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		admin.Username,
		admin.Email,
		admin.Role,
		admin.IsActive,
		admin.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AdminRepository) DeleteAdminUser(adminID int) error {
	query := `DELETE FROM admin_users WHERE id = $1`
	_, err := r.db.Exec(query, adminID)
	return err
}

func (r *AdminRepository) GetAdminUserByID(adminID int) (*models.AdminUser, error) {
	query := `
		SELECT id, username, email, role, is_active, last_login, created_at
		FROM admin_users
		WHERE id = $1
	`

	var admin models.AdminUser
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, adminID).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.Role,
		&admin.IsActive,
		&lastLogin,
		&admin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin by ID: %w", err)
	}

	if lastLogin.Valid {
		admin.LastLogin = &lastLogin.Time
	}

	return &admin, nil
}

func (r *AdminRepository) UpdateAdminPassword(adminID int, passwordHash string) error {
	query := `UPDATE admin_users SET password_hash = $1 WHERE id = $2`
	result, err := r.db.Exec(query, passwordHash, adminID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AdminRepository) GetTodayStats() (*models.AdminStats, error) {
	today := time.Now().Truncate(24 * time.Hour)
	startOfDay := today
	endOfDay := today.Add(24 * time.Hour)

	return r.GetDashboardStats(startOfDay, endOfDay)
}

func (r *AdminRepository) GetAdminUsersCount() (int, error) {
	query := `SELECT COUNT(*) FROM admin_users`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get admin users count: %w", err)
	}
	return count, nil
}
