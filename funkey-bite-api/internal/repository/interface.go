package repository

import (
	"database/sql"

	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type IAdminRepository interface {
	GetDashboardStats(fromDate, toDate time.Time) (*models.AdminStats, error)
	GetSalesReport(fromDate, toDate time.Time) ([]models.SalesReport, error)
	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) error
	DeleteMenuItem(id int) error
	GetAllOrders(limit, offset int, status string) ([]models.Order, error)
	GetAllUsers(limit, offset int) ([]models.User, error)
	UpdateUserStatus(userID int, isActive bool) error
	GetOrdersCount(status string) (int, error) // Add this
	GetUsersCount() (int, error)               // Add this
}

// Add these to IAdminRepository or create new interface
type IOrderRepository interface {
	Create(order *models.Order) (*models.Order, error)
	CreateOrderItem(item *models.OrderItem) (*models.OrderItem, error)
	GetOrderWithItems(id int) (*models.Order, error)
	GetOrdersByUserID(userID int) ([]models.Order, error)
	UpdateOrderStatus(id int, status string) error
	GetOrderByOrderNumber(orderNumber string) (*models.Order, error)
	GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error)
	CancelOrder(id int) error

	// Transaction methods
	BeginTransaction() (*sql.Tx, error)
	CreateOrderWithTransaction(tx *sql.Tx, order *models.Order) (*models.Order, error)
	CreateOrderItemWithTransaction(tx *sql.Tx, item *models.OrderItem) (*models.OrderItem, error)
}
