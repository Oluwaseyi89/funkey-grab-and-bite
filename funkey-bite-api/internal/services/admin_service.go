package services

import (
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type AdminService interface {
	GetDashboardStats() (*models.AdminStats, error)
	GetSalesReport(fromDate, toDate string) ([]models.SalesReport, error)
	GetAllOrders(page, limit int, status string) ([]models.Order, int, error)
	UpdateOrderStatus(orderID int, status string) error
	GetAllUsers(page, limit int) ([]models.User, int, error)
	UpdateUserStatus(userID int, isActive bool) error
	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) error
	DeleteMenuItem(id int) error
	GetMenuItemByID(id int) (*models.MenuItem, error) // Add this
	GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error)
}

type adminService struct {
	adminRepo    repository.IAdminRepository // Changed to IAdminRepository
	orderRepo    repository.OrderRepository
	userRepo     repository.UserRepository
	cateringRepo repository.CateringRepository
	menuRepo     repository.MenuRepository
}

func NewAdminService(
	adminRepo repository.IAdminRepository, // Changed to IAdminRepository
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	cateringRepo repository.CateringRepository,
	menuRepo repository.MenuRepository,
) AdminService {
	return &adminService{
		adminRepo:    adminRepo,
		orderRepo:    orderRepo,
		userRepo:     userRepo,
		cateringRepo: cateringRepo,
		menuRepo:     menuRepo,
	}
}

func (s *adminService) GetDashboardStats() (*models.AdminStats, error) {
	// Last 30 days by default
	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	return s.adminRepo.GetDashboardStats(fromDate, toDate)
}

func (s *adminService) GetSalesReport(fromDateStr, toDateStr string) ([]models.SalesReport, error) {
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid from date format")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid to date format")
	}

	if fromDate.After(toDate) {
		return nil, fmt.Errorf("from date cannot be after to date")
	}

	return s.adminRepo.GetSalesReport(fromDate, toDate)
}

func (s *adminService) GetAllOrders(page, limit int, status string) ([]models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	orders, err := s.adminRepo.GetAllOrders(limit, offset, status)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination using the new interface method
	totalCount, err := s.adminRepo.GetOrdersCount(status)
	if err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (s *adminService) UpdateOrderStatus(orderID int, status string) error {
	return s.orderRepo.UpdateOrderStatus(orderID, status)
}

func (s *adminService) GetAllUsers(page, limit int) ([]models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	users, err := s.adminRepo.GetAllUsers(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count using the new interface method
	totalCount, err := s.adminRepo.GetUsersCount()
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (s *adminService) UpdateUserStatus(userID int, isActive bool) error {
	return s.adminRepo.UpdateUserStatus(userID, isActive)
}

func (s *adminService) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	// Validate category exists
	categories, err := s.menuRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	categoryExists := false
	for _, cat := range categories {
		if cat.ID == item.CategoryID {
			categoryExists = true
			break
		}
	}

	if !categoryExists {
		return nil, fmt.Errorf("category does not exist")
	}

	return s.adminRepo.CreateMenuItem(item)
}

func (s *adminService) UpdateMenuItem(item *models.MenuItem) error {
	// Check if item exists
	existing, err := s.menuRepo.GetByID(item.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("menu item not found")
	}

	return s.adminRepo.UpdateMenuItem(item)
}

func (s *adminService) DeleteMenuItem(id int) error {
	// Check if item exists
	existing, err := s.menuRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("menu item not found")
	}

	return s.adminRepo.DeleteMenuItem(id)
}

func (s *adminService) GetMenuItemByID(id int) (*models.MenuItem, error) {
	return s.menuRepo.GetByID(id)
}

func (s *adminService) GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error) {
	// TODO: Implement when catering repository has pagination methods
	// For now, return empty
	return []models.CateringRequest{}, 0, nil
}

// package services

// import (
// 	"fmt"
// 	"time"

// 	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
// 	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
// )

// type AdminService interface {
// 	GetDashboardStats() (*models.AdminStats, error)
// 	GetSalesReport(fromDate, toDate string) ([]models.SalesReport, error)
// 	GetAllOrders(page, limit int, status string) ([]models.Order, int, error)
// 	UpdateOrderStatus(orderID int, status string) error
// 	GetAllUsers(page, limit int) ([]models.User, int, error)
// 	UpdateUserStatus(userID int, isActive bool) error
// 	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
// 	UpdateMenuItem(item *models.MenuItem) error
// 	DeleteMenuItem(id int) error
// 	GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error)
// }

// type adminService struct {
// 	adminRepo    repository.AdminRepository
// 	orderRepo    repository.OrderRepository
// 	userRepo     repository.UserRepository
// 	cateringRepo repository.CateringRepository
// 	menuRepo     repository.MenuRepository
// }

// func NewAdminService(
// 	adminRepo repository.AdminRepository,
// 	orderRepo repository.OrderRepository,
// 	userRepo repository.UserRepository,
// 	cateringRepo repository.CateringRepository,
// 	menuRepo repository.MenuRepository,
// ) AdminService {
// 	return &adminService{
// 		adminRepo:    adminRepo,
// 		orderRepo:    orderRepo,
// 		userRepo:     userRepo,
// 		cateringRepo: cateringRepo,
// 		menuRepo:     menuRepo,
// 	}
// }

// func (s *adminService) GetDashboardStats() (*models.AdminStats, error) {
// 	// Last 30 days by default
// 	fromDate := time.Now().AddDate(0, 0, -30)
// 	toDate := time.Now()

// 	return s.adminRepo.GetDashboardStats(fromDate, toDate)
// }

// func (s *adminService) GetSalesReport(fromDateStr, toDateStr string) ([]models.SalesReport, error) {
// 	fromDate, err := time.Parse("2006-01-02", fromDateStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid from date format")
// 	}

// 	toDate, err := time.Parse("2006-01-02", toDateStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid to date format")
// 	}

// 	if fromDate.After(toDate) {
// 		return nil, fmt.Errorf("from date cannot be after to date")
// 	}

// 	return s.adminRepo.GetSalesReport(fromDate, toDate)
// }

// func (s *adminService) GetAllOrders(page, limit int, status string) ([]models.Order, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	orders, err := s.adminRepo.GetAllOrders(limit, offset, status)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Get total count for pagination
// 	var totalCount int
// 	query := `SELECT COUNT(*) FROM orders WHERE ($1 = '' OR status = $1)`
// 	err = s.adminRepo.db.QueryRow(query, status).Scan(&totalCount)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return orders, totalCount, nil
// }

// func (s *adminService) UpdateOrderStatus(orderID int, status string) error {
// 	return s.orderRepo.UpdateOrderStatus(orderID, status)
// }

// func (s *adminService) GetAllUsers(page, limit int) ([]models.User, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	users, err := s.adminRepo.GetAllUsers(limit, offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Get total count
// 	var totalCount int
// 	err = s.adminRepo.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalCount)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return users, totalCount, nil
// }

// func (s *adminService) UpdateUserStatus(userID int, isActive bool) error {
// 	return s.adminRepo.UpdateUserStatus(userID, isActive)
// }

// func (s *adminService) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
// 	// Validate category exists
// 	categories, err := s.menuRepo.GetCategories()
// 	if err != nil {
// 		return nil, err
// 	}

// 	categoryExists := false
// 	for _, cat := range categories {
// 		if cat.ID == item.CategoryID {
// 			categoryExists = true
// 			break
// 		}
// 	}

// 	if !categoryExists {
// 		return nil, fmt.Errorf("category does not exist")
// 	}

// 	return s.adminRepo.CreateMenuItem(item)
// }

// func (s *adminService) UpdateMenuItem(item *models.MenuItem) error {
// 	// Check if item exists
// 	existing, err := s.menuRepo.GetByID(item.ID)
// 	if err != nil {
// 		return err
// 	}
// 	if existing == nil {
// 		return fmt.Errorf("menu item not found")
// 	}

// 	return s.adminRepo.UpdateMenuItem(item)
// }

// func (s *adminService) DeleteMenuItem(id int) error {
// 	// Check if item exists
// 	existing, err := s.menuRepo.GetByID(id)
// 	if err != nil {
// 		return err
// 	}
// 	if existing == nil {
// 		return fmt.Errorf("menu item not found")
// 	}

// 	return s.adminRepo.DeleteMenuItem(id)
// }

// func (s *adminService) GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error) {
// 	// Implementation would go here
// 	// For now, return empty
// 	return []models.CateringRequest{}, 0, nil
// }
