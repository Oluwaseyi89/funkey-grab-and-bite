package services

import (
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type UserService interface {
	GetByID(id int) (*models.User, error)
	UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error)
	GetOrderHistory(userID int) ([]models.Order, error)
	ChangePassword(userID int, currentPassword, newPassword string) error
}

type userService struct {
	userRepo  repository.UserRepository
	orderRepo repository.OrderRepository
}

func NewUserService(userRepo repository.UserRepository, orderRepo repository.OrderRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		orderRepo: orderRepo,
	}
}

func (s *userService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, nil // User not found
	}

	// Remove sensitive information before returning
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error) {
	// First, get the current user to validate
	currentUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Prepare updates map for repository
	updateFields := make(map[string]interface{})

	// Only include fields that are provided and different from current
	if updates.FullName != nil && *updates.FullName != currentUser.FullName {
		updateFields["full_name"] = *updates.FullName
	}

	if updates.Email != nil && *updates.Email != "" {
		// Check if email is already used by another user
		existingUser, err := s.userRepo.FindByPhoneOrEmail("", *updates.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email availability: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("email already in use")
		}
		updateFields["email"] = *updates.Email
	}

	if updates.Phone != nil && *updates.Phone != currentUser.Phone {
		// Check if phone is already used by another user
		existingUser, err := s.userRepo.FindByPhone(*updates.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to check phone availability: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("phone number already in use")
		}
		updateFields["phone"] = *updates.Phone
	}

	// If no fields to update, return current user
	if len(updateFields) == 0 {
		currentUser.PasswordHash = "" // Remove sensitive data
		return currentUser, nil
	}

	// Apply updates
	err = s.userRepo.UpdateProfile(id, updateFields)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Return updated user
	return s.GetByID(id)
}

func (s *userService) GetOrderHistory(userID int) ([]models.Order, error) {
	orders, err := s.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}
	return orders, nil
}

func (s *userService) ChangePassword(userID int, currentPassword, newPassword string) error {
	// Get user with password hash
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password strength
	if !utils.ValidatePasswordStrength(newPassword) {
		return fmt.Errorf("new password must be at least 8 characters with uppercase, lowercase, and numbers")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password in repository
	err = s.userRepo.UpdateProfile(userID, map[string]interface{}{
		"password_hash": hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// package services

// import (
// 	"fmt"

// 	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
// 	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
// 	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
// )

// type UserService interface {
// 	GetByID(id int) (*models.User, error)
// 	UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error)
// 	GetOrderHistory(userID int, page, limit int) ([]models.Order, int, error) // Updated with pagination
// 	ChangePassword(userID int, currentPassword, newPassword string) error
// 	GetCateringRequests(userID int, page, limit int) ([]models.CateringRequest, int, error) // Add this
// }

// type userService struct {
// 	userRepo     repository.UserRepository
// 	orderRepo    repository.OrderRepository
// 	cateringRepo repository.CateringRepository // Add this
// }

// func NewUserService(userRepo repository.UserRepository, orderRepo repository.OrderRepository, cateringRepo repository.CateringRepository) UserService {
// 	return &userService{
// 		userRepo:     userRepo,
// 		orderRepo:    orderRepo,
// 		cateringRepo: cateringRepo, // Add this
// 	}
// }

// func (s *userService) GetByID(id int) (*models.User, error) {
// 	user, err := s.userRepo.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user: %w", err)
// 	}
// 	if user == nil {
// 		return nil, nil // User not found
// 	}

// 	// Remove sensitive information before returning
// 	user.PasswordHash = ""
// 	return user, nil
// }

// func (s *userService) UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error) {
// 	// First, get the current user to validate
// 	currentUser, err := s.userRepo.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get current user: %w", err)
// 	}
// 	if currentUser == nil {
// 		return nil, fmt.Errorf("user not found")
// 	}

// 	// Prepare updates map for repository
// 	updateFields := make(map[string]interface{})

// 	// Only include fields that are provided and different from current
// 	if updates.FullName != nil && *updates.FullName != currentUser.FullName {
// 		updateFields["full_name"] = *updates.FullName
// 	}

// 	if updates.Email != nil && *updates.Email != "" {
// 		// Check if email is already used by another user
// 		existingUser, err := s.userRepo.FindByPhoneOrEmail("", *updates.Email)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to check email availability: %w", err)
// 		}
// 		if existingUser != nil && existingUser.ID != id {
// 			return nil, fmt.Errorf("email already in use")
// 		}
// 		updateFields["email"] = *updates.Email
// 	}

// 	if updates.Phone != nil && *updates.Phone != currentUser.Phone {
// 		// Check if phone is already used by another user
// 		existingUser, err := s.userRepo.FindByPhone(*updates.Phone)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to check phone availability: %w", err)
// 		}
// 		if existingUser != nil && existingUser.ID != id {
// 			return nil, fmt.Errorf("phone number already in use")
// 		}
// 		updateFields["phone"] = *updates.Phone
// 	}

// 	// If no fields to update, return current user
// 	if len(updateFields) == 0 {
// 		currentUser.PasswordHash = "" // Remove sensitive data
// 		return currentUser, nil
// 	}

// 	// Apply updates
// 	err = s.userRepo.UpdateProfile(id, updateFields)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update profile: %w", err)
// 	}

// 	// Return updated user
// 	return s.GetByID(id)
// }

// func (s *userService) GetOrderHistory(userID int, page, limit int) ([]models.Order, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	// Get orders with pagination
// 	orders, err := s.orderRepo.GetOrdersByUserIDWithPagination(userID, limit, offset)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get order history: %w", err)
// 	}

// 	// Get total count
// 	totalCount, err := s.orderRepo.GetOrdersCountByUserID(userID)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get orders count: %w", err)
// 	}

// 	return orders, totalCount, nil
// }

// // Keep backward compatibility method
// func (s *userService) GetOrderHistorySimple(userID int) ([]models.Order, error) {
// 	orders, err := s.orderRepo.GetOrdersByUserID(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get order history: %w", err)
// 	}
// 	return orders, nil
// }

// func (s *userService) GetCateringRequests(userID int, page, limit int) ([]models.CateringRequest, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	// Get catering requests with pagination
// 	requests, err := s.cateringRepo.GetByUserIDWithPagination(userID, limit, offset)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get catering requests: %w", err)
// 	}

// 	// Get total count
// 	totalCount, err := s.cateringRepo.GetCountByUserID(userID)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get catering requests count: %w", err)
// 	}

// 	return requests, totalCount, nil
// }

// // Keep backward compatibility method
// func (s *userService) GetCateringRequestsSimple(userID int) ([]models.CateringRequest, error) {
// 	requests, err := s.cateringRepo.GetByUserID(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get catering requests: %w", err)
// 	}
// 	return requests, nil
// }

// func (s *userService) ChangePassword(userID int, currentPassword, newPassword string) error {
// 	// Get user with password hash
// 	user, err := s.userRepo.GetByID(userID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get user: %w", err)
// 	}
// 	if user == nil {
// 		return fmt.Errorf("user not found")
// 	}

// 	// Verify current password
// 	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
// 		return fmt.Errorf("current password is incorrect")
// 	}

// 	// Validate new password strength
// 	if !utils.ValidatePasswordStrength(newPassword) {
// 		return fmt.Errorf("new password must be at least 8 characters with uppercase, lowercase, and numbers")
// 	}

// 	// Hash new password
// 	hashedPassword, err := utils.HashPassword(newPassword)
// 	if err != nil {
// 		return fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	// Update password in repository
// 	err = s.userRepo.UpdateProfile(userID, map[string]interface{}{
// 		"password_hash": hashedPassword,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to update password: %w", err)
// 	}

// 	return nil
// }

// package services

// import (
// 	"fmt"

// 	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
// 	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
// 	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
// )

// type UserService interface {
// 	GetByID(id int) (*models.User, error)
// 	UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error)
// 	GetOrderHistory(userID int, page, limit int) ([]models.Order, int, error) // Updated with pagination
// 	ChangePassword(userID int, currentPassword, newPassword string) error
// 	GetCateringRequests(userID int, page, limit int) ([]models.CateringRequest, int, error) // Add this
// }

// type userService struct {
// 	userRepo     repository.UserRepository
// 	orderRepo    repository.OrderRepository
// 	cateringRepo repository.CateringRepository // Add this
// }

// func NewUserService(userRepo repository.UserRepository, orderRepo repository.OrderRepository, cateringRepo repository.CateringRepository) UserService {
// 	return &userService{
// 		userRepo:     userRepo,
// 		orderRepo:    orderRepo,
// 		cateringRepo: cateringRepo, // Add this
// 	}
// }

// func (s *userService) GetByID(id int) (*models.User, error) {
// 	user, err := s.userRepo.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user: %w", err)
// 	}
// 	if user == nil {
// 		return nil, nil // User not found
// 	}

// 	// Remove sensitive information before returning
// 	user.PasswordHash = ""
// 	return user, nil
// }

// func (s *userService) UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error) {
// 	// First, get the current user to validate
// 	currentUser, err := s.userRepo.GetByID(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get current user: %w", err)
// 	}
// 	if currentUser == nil {
// 		return nil, fmt.Errorf("user not found")
// 	}

// 	// Prepare updates map for repository
// 	updateFields := make(map[string]interface{})

// 	// Only include fields that are provided and different from current
// 	if updates.FullName != nil && *updates.FullName != currentUser.FullName {
// 		updateFields["full_name"] = *updates.FullName
// 	}

// 	if updates.Email != nil && *updates.Email != "" {
// 		// Check if email is already used by another user
// 		existingUser, err := s.userRepo.FindByPhoneOrEmail("", *updates.Email)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to check email availability: %w", err)
// 		}
// 		if existingUser != nil && existingUser.ID != id {
// 			return nil, fmt.Errorf("email already in use")
// 		}
// 		updateFields["email"] = *updates.Email
// 	}

// 	if updates.Phone != nil && *updates.Phone != currentUser.Phone {
// 		// Check if phone is already used by another user
// 		existingUser, err := s.userRepo.FindByPhone(*updates.Phone)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to check phone availability: %w", err)
// 		}
// 		if existingUser != nil && existingUser.ID != id {
// 			return nil, fmt.Errorf("phone number already in use")
// 		}
// 		updateFields["phone"] = *updates.Phone
// 	}

// 	// If no fields to update, return current user
// 	if len(updateFields) == 0 {
// 		currentUser.PasswordHash = "" // Remove sensitive data
// 		return currentUser, nil
// 	}

// 	// Apply updates
// 	err = s.userRepo.UpdateProfile(id, updateFields)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update profile: %w", err)
// 	}

// 	// Return updated user
// 	return s.GetByID(id)
// }

// func (s *userService) GetOrderHistory(userID int, page, limit int) ([]models.Order, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	// Get orders with pagination
// 	orders, err := s.orderRepo.GetOrdersByUserIDWithPagination(userID, limit, offset)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get order history: %w", err)
// 	}

// 	// Get total count
// 	totalCount, err := s.orderRepo.GetOrdersCountByUserID(userID)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get orders count: %w", err)
// 	}

// 	return orders, totalCount, nil
// }

// // Keep backward compatibility method
// func (s *userService) GetOrderHistorySimple(userID int) ([]models.Order, error) {
// 	orders, err := s.orderRepo.GetOrdersByUserID(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get order history: %w", err)
// 	}
// 	return orders, nil
// }

// func (s *userService) GetCateringRequests(userID int, page, limit int) ([]models.CateringRequest, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	offset := (page - 1) * limit

// 	// Get catering requests with pagination
// 	requests, err := s.cateringRepo.GetByUserIDWithPagination(userID, limit, offset)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get catering requests: %w", err)
// 	}

// 	// Get total count
// 	totalCount, err := s.cateringRepo.GetCountByUserID(userID)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get catering requests count: %w", err)
// 	}

// 	return requests, totalCount, nil
// }

// // Keep backward compatibility method
// func (s *userService) GetCateringRequestsSimple(userID int) ([]models.CateringRequest, error) {
// 	requests, err := s.cateringRepo.GetByUserID(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get catering requests: %w", err)
// 	}
// 	return requests, nil
// }

// func (s *userService) ChangePassword(userID int, currentPassword, newPassword string) error {
// 	// Get user with password hash
// 	user, err := s.userRepo.GetByID(userID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get user: %w", err)
// 	}
// 	if user == nil {
// 		return fmt.Errorf("user not found")
// 	}

// 	// Verify current password
// 	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
// 		return fmt.Errorf("current password is incorrect")
// 	}

// 	// Validate new password strength
// 	if !utils.ValidatePasswordStrength(newPassword) {
// 		return fmt.Errorf("new password must be at least 8 characters with uppercase, lowercase, and numbers")
// 	}

// 	// Hash new password
// 	hashedPassword, err := utils.HashPassword(newPassword)
// 	if err != nil {
// 		return fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	// Update password in repository
// 	err = s.userRepo.UpdateProfile(userID, map[string]interface{}{
// 		"password_hash": hashedPassword,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to update password: %w", err)
// 	}

// 	return nil
// }////s.orderRepo.GetOrdersByUserIDWithPagination undefined (type repository.OrderRepository has no field or method GetOrdersByUserIDWithPagination)compilerMissingFieldOrMethod
// ////s.orderRepo.GetOrdersCountByUserID undefined (type repository.OrderRepository has no field or method GetOrdersCountByUserID)compilerMissingFieldOrMethod
// ////s.cateringRepo.GetByUserIDWithPagination undefined (type repository.CateringRepository has no field or method GetByUserIDWithPagination)compilerMissingFieldOrMethod
// ////s.cateringRepo.GetCountByUserID undefined (type repository.CateringRepository has no field or method GetCountByUserID)compilerMissingFieldOrMethod
