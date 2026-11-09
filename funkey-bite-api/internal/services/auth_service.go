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
	FindOrderCustomer(phone string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(userData models.UserRegistration) (*models.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByPhoneOrEmail(userData.Phone, userData.Email)
	isGuest := err == nil && existingUser != nil && existingUser.PasswordHash == "" && existingUser.Phone == userData.Phone
	if err == nil && existingUser != nil && !isGuest {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		return nil, err
	}

	if isGuest {
		return s.claimGuestAccount(existingUser, userData, hashedPassword)
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
	if err != nil || user == nil {
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

// FindOrderCustomer returns the existing customer record for a phone number,
// or nil when there is none. It never creates one.
func (s *authService) FindOrderCustomer(phone string) (*models.User, error) {
	return s.userRepo.FindByPhone(phone)
}

// claimGuestAccount turns a guest customer record into a full account.
func (s *authService) claimGuestAccount(user *models.User, userData models.UserRegistration, hashedPassword string) (*models.AuthResponse, error) {
	if err := s.userRepo.SetPassword(user.ID, hashedPassword); err != nil {
		return nil, err
	}
	if err := s.userRepo.UpdateProfile(user.ID, map[string]interface{}{
		"full_name": userData.FullName,
		"email":     userData.Email,
	}); err != nil {
		return nil, err
	}

	user.FullName = userData.FullName
	user.Email = &userData.Email
	user.PasswordHash = hashedPassword

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
