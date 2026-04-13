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
		return nil, nil
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *userService) UpdateProfile(id int, updates *models.ProfileUpdate) (*models.User, error) {
	currentUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	updateFields := make(map[string]interface{})

	if updates.FullName != nil && *updates.FullName != currentUser.FullName {
		updateFields["full_name"] = *updates.FullName
	}

	if updates.Email != nil && *updates.Email != "" {
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
		existingUser, err := s.userRepo.FindByPhone(*updates.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to check phone availability: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("phone number already in use")
		}
		updateFields["phone"] = *updates.Phone
	}

	if len(updateFields) == 0 {
		currentUser.PasswordHash = ""
		return currentUser, nil
	}

	err = s.userRepo.UpdateProfile(id, updateFields)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

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
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	if !utils.ValidatePasswordStrength(newPassword) {
		return fmt.Errorf("new password must be at least 8 characters with uppercase, lowercase, and numbers")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.userRepo.UpdateProfile(userID, map[string]interface{}{
		"password_hash": hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
