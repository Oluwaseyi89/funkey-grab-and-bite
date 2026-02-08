package repository

import (
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
