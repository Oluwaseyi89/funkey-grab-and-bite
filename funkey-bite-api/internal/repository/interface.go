package repository

import (
	"database/sql"

	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type IAdminRepository interface {
	GetDashboardStats(fromDate, toDate time.Time) (*models.AdminStats, error)
	GetTodayStats() (*models.AdminStats, error)
	GetSalesReport(fromDate, toDate time.Time) ([]models.SalesReport, error)

	GetAdminByEmail(email string) (*models.AdminUser, error)
	UpdateAdminLastLogin(adminID int) error
	GetAdminUsers(limit, offset int) ([]models.AdminUser, error)
	GetAdminUsersCount() (int, error)
	GetAdminUserByID(adminID int) (*models.AdminUser, error)
	CreateAdminUser(admin *models.AdminUser) error
	UpdateAdminUser(admin *models.AdminUser) error
	DeleteAdminUser(adminID int) error
	UpdateAdminPassword(adminID int, passwordHash string) error

	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) error
	DeleteMenuItem(id int) error

	GetAllOrders(limit, offset int, filter models.OrderListFilter) ([]models.Order, error)
	GetOrdersCount(filter models.OrderListFilter) (int, error)

	GetAllUsers(limit, offset int, filter models.CustomerListFilter) ([]models.AdminCustomer, error)
	GetUsersCount(filter models.CustomerListFilter) (int, error)
	GetCustomer(userID int) (*models.AdminCustomer, error)
	UpdateUserStatus(userID int, isActive bool) error
}

type IOrderRepository interface {
	Create(order *models.Order) (*models.Order, error)
	CreateOrderItem(item *models.OrderItem) (*models.OrderItem, error)
	GetOrderWithItems(id int) (*models.Order, error)
	GetOrdersByUserID(userID int) ([]models.Order, error)
	UpdateOrderStatus(id int, status string) error
	GetOrderByOrderNumber(orderNumber string) (*models.Order, error)
	GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error)
	CancelOrder(id int) error
	SetOrderCustomer(orderID, userID int) error

	BeginTransaction() (*sql.Tx, error)
	CreateOrderWithTransaction(tx *sql.Tx, order *models.Order) (*models.Order, error)
	CreateOrderItemWithTransaction(tx *sql.Tx, item *models.OrderItem) (*models.OrderItem, error)

	GetOrderByPaymentReference(reference string) (*models.Order, error)
	SetOrderPaymentInitiated(orderID int, reference, accountNumber, accountName, bankName string, expiresAt *time.Time) error
	MarkOrderPaymentPaid(orderID int, paidAt time.Time) error
	MarkPaymentCollected(orderID int, via string, paidAt time.Time) (bool, error)
}
