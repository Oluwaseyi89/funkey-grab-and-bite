package services

import (
	"errors"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type AuthService interface {
	Register(userData models.UserRegistration) (*models.AuthResponse, error)
	Login(phone, password string) (*models.AuthResponse, error)
	CheckUserExists(phone, email string) (*models.User, bool, error)
	AuthenticateOrder(orderData models.OrderWithAuth) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(userData models.UserRegistration) (*models.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByPhoneOrEmail(userData.Phone, userData.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Phone:        userData.Phone,
		Email:        &userData.Email,
		FullName:     userData.FullName,
		PasswordHash: hashedPassword,
		IsVerified:   true,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	user, err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:        user,
		AccessToken: token,
		IsNewUser:   true,
	}, nil
}

func (s *authService) Login(phone, password string) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLogin = &now
	s.userRepo.UpdateLastLogin(user.ID, now)

	token, err := utils.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:        user,
		AccessToken: token,
	}, nil
}

func (s *authService) CheckUserExists(phone, email string) (*models.User, bool, error) {
	user, err := s.userRepo.FindByPhoneOrEmail(phone, email)
	if err != nil || user == nil {
		return nil, false, nil
	}

	return user, true, nil
}

func (s *authService) AuthenticateOrder(orderData models.OrderWithAuth) (*models.User, error) {
	user, exists, err := s.CheckUserExists(orderData.CustomerPhone, "")
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("user_not_found")
	}

	if orderData.Password == nil || *orderData.Password == "" {
		return nil, errors.New("password_required")
	}

	if !utils.CheckPasswordHash(*orderData.Password, user.PasswordHash) {
		return nil, errors.New("invalid_password")
	}

	return user, nil
}
