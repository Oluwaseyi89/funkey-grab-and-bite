package services

import (
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// Update the AdminService interface to include admin auth methods
type AdminService interface {
	// Admin Authentication
	AdminLogin(email, password string) (*models.AdminUser, string, error)
	AdminLogout(adminID int) error

	// Admin User Management
	GetAdminUsers(page, limit int) ([]models.AdminUser, int, error)
	GetAdminUserByID(adminID int) (*models.AdminUser, error)
	CreateAdminUser(admin *models.AdminUser, password string) (*models.AdminUser, error)
	UpdateAdminUser(adminID int, updates *models.AdminUser) error
	DeleteAdminUser(adminID int) error
	UpdateAdminPassword(adminID int, currentPassword, newPassword string) error

	// Dashboard
	GetDashboardStats() (*models.AdminStats, error)
	GetTodayStats() (*models.AdminStats, error)
	GetSalesReport(fromDate, toDate string) ([]models.SalesReport, error)

	// Orders Management
	GetAllOrders(page, limit int, status string) ([]models.Order, int, error)
	UpdateOrderStatus(orderID int, status string) error

	// User Management
	GetAllUsers(page, limit int) ([]models.User, int, error)
	UpdateUserStatus(userID int, isActive bool) error

	// Menu Management
	CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error)
	UpdateMenuItem(item *models.MenuItem) error
	DeleteMenuItem(id int) error
	GetMenuItemByID(id int) (*models.MenuItem, error)

	// Catering Management
	GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error)
}

type adminService struct {
	adminRepo    repository.IAdminRepository
	orderRepo    repository.OrderRepository
	userRepo     repository.UserRepository
	cateringRepo repository.CateringRepository
	menuRepo     repository.MenuRepository
}

func NewAdminService(
	adminRepo repository.IAdminRepository,
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

// Admin Authentication Methods

func (s *adminService) AdminLogin(email, password string) (*models.AdminUser, string, error) {
	// Get admin by email
	admin, err := s.adminRepo.GetAdminByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get admin: %w", err)
	}
	if admin == nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check if admin is active
	if !admin.IsActive {
		return nil, "", fmt.Errorf("admin account is inactive")
	}

	// Verify password
	if !utils.VerifyPassword(password, admin.PasswordHash) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateAdminJWT(admin.ID, admin.Email, admin.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	err = s.adminRepo.UpdateAdminLastLogin(admin.ID)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	return admin, token, nil
}

func (s *adminService) AdminLogout(adminID int) error {
	// For JWT, logout is client-side (token invalidation)
	// Could implement token blacklist if needed
	return nil
}

// Admin User Management Methods

// func (s *adminService) GetAdminUsers(page, limit int) ([]models.AdminUser, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	admins, err := s.adminRepo.GetAdminUsers(limit, offset)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get admin users: %w", err)
// 	}

// 	// TODO: Need to add GetAdminUsersCount method to repository
// 	// For now, return length as total
// 	total := len(admins)

// 	return admins, total, nil
// }

func (s *adminService) GetAdminUsers(page, limit int) ([]models.AdminUser, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	admins, err := s.adminRepo.GetAdminUsers(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get admin users: %w", err)
	}

	// Get total count using the new method
	totalCount, err := s.adminRepo.GetAdminUsersCount()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get admin users count: %w", err)
	}

	return admins, totalCount, nil
}

func (s *adminService) GetAdminUserByID(adminID int) (*models.AdminUser, error) {
	return s.adminRepo.GetAdminUserByID(adminID)
}

func (s *adminService) CreateAdminUser(admin *models.AdminUser, password string) (*models.AdminUser, error) {
	// Validate password
	if !utils.ValidatePasswordStrength(password) {
		return nil, fmt.Errorf("password must be at least 8 characters with uppercase, lowercase, and numbers")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	admin.PasswordHash = hashedPassword

	// Set default values if not provided
	if admin.Role == "" {
		admin.Role = "manager"
	}

	err = s.adminRepo.CreateAdminUser(admin)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	// Clear password hash from response
	admin.PasswordHash = ""

	return admin, nil
}

func (s *adminService) UpdateAdminUser(adminID int, updates *models.AdminUser) error {
	// Get existing admin
	existing, err := s.adminRepo.GetAdminUserByID(adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("admin not found")
	}

	// Update fields
	if updates.Username != "" {
		existing.Username = updates.Username
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.Role != "" {
		existing.Role = updates.Role
	}

	return s.adminRepo.UpdateAdminUser(existing)
}

func (s *adminService) DeleteAdminUser(adminID int) error {
	// Prevent deleting own account (check in handler)
	return s.adminRepo.DeleteAdminUser(adminID)
}

func (s *adminService) UpdateAdminPassword(adminID int, currentPassword, newPassword string) error {
	// Get admin
	admin, err := s.adminRepo.GetAdminUserByID(adminID)
	if err != nil {
		return fmt.Errorf("failed to get admin: %w", err)
	}
	if admin == nil {
		return fmt.Errorf("admin not found")
	}

	// Verify current password
	if !utils.VerifyPassword(currentPassword, admin.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if !utils.ValidatePasswordStrength(newPassword) {
		return fmt.Errorf("new password must be at least 8 characters with uppercase, lowercase, and numbers")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.adminRepo.UpdateAdminPassword(adminID, hashedPassword)
}

// Dashboard Methods

func (s *adminService) GetDashboardStats() (*models.AdminStats, error) {
	// Last 30 days by default
	fromDate := time.Now().AddDate(0, 0, -30)
	toDate := time.Now()

	return s.adminRepo.GetDashboardStats(fromDate, toDate)
}

func (s *adminService) GetTodayStats() (*models.AdminStats, error) {
	return s.adminRepo.GetTodayStats()
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

// Order Management Methods (unchanged)

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

	totalCount, err := s.adminRepo.GetOrdersCount(status)
	if err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (s *adminService) UpdateOrderStatus(orderID int, status string) error {
	return s.orderRepo.UpdateOrderStatus(orderID, status)
}

// User Management Methods (unchanged)

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

	totalCount, err := s.adminRepo.GetUsersCount()
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (s *adminService) UpdateUserStatus(userID int, isActive bool) error {
	return s.adminRepo.UpdateUserStatus(userID, isActive)
}

// Menu Management Methods (unchanged)

func (s *adminService) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
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

// Catering Management Methods (unchanged)

func (s *adminService) GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	requests, err := s.cateringRepo.GetAll()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get catering requests: %w", err)
	}

	if status != "" {
		filtered := []models.CateringRequest{}
		for _, req := range requests {
			if string(req.Status) == status {
				filtered = append(filtered, req)
			}
		}
		requests = filtered
	}

	total := len(requests)

	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	if start >= total {
		return []models.CateringRequest{}, total, nil
	}

	return requests[start:end], total, nil
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
// 	GetMenuItemByID(id int) (*models.MenuItem, error) // Add this
// 	GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error)
// }

// type adminService struct {
// 	adminRepo    repository.IAdminRepository // Changed to IAdminRepository
// 	orderRepo    repository.OrderRepository
// 	userRepo     repository.UserRepository
// 	cateringRepo repository.CateringRepository
// 	menuRepo     repository.MenuRepository
// }

// func NewAdminService(
// 	adminRepo repository.IAdminRepository, // Changed to IAdminRepository
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

// 	// Get total count for pagination using the new interface method
// 	totalCount, err := s.adminRepo.GetOrdersCount(status)
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

// 	// Get total count using the new interface method
// 	totalCount, err := s.adminRepo.GetUsersCount()
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

// func (s *adminService) GetMenuItemByID(id int) (*models.MenuItem, error) {
// 	return s.menuRepo.GetByID(id)
// }

// func (s *adminService) GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	// TODO: Need to add these methods to CateringRepository interface
// 	// For now, use GetAll() as fallback
// 	requests, err := s.cateringRepo.GetAll()
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get catering requests: %w", err)
// 	}

// 	// Simple filtering by status
// 	if status != "" {
// 		filtered := []models.CateringRequest{}
// 		for _, req := range requests {
// 			if string(req.Status) == status {
// 				filtered = append(filtered, req)
// 			}
// 		}
// 		requests = filtered
// 	}

// 	total := len(requests)

// 	// Simple pagination
// 	start := offset
// 	if start > total {
// 		start = total
// 	}
// 	end := start + limit
// 	if end > total {
// 		end = total
// 	}

// 	if start >= total {
// 		return []models.CateringRequest{}, total, nil
// 	}

// 	return requests[start:end], total, nil
// }
