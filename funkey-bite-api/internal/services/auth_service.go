package services

import (
	"errors"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type AuthService interface {
	Login(phone, password string) (*models.AuthResponse, error)
	CheckUserExists(phone, email string) (*models.User, bool, error)
	FindOrderCustomer(phone, email string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(phone, password string) (*models.AuthResponse, error) {
	// Numbers match in any format; a guest record has no password, so it
	// can't be logged into.
	user, err := s.userRepo.FindByPhoneAnyFormat(utils.CleanPhone(phone))
	if err != nil || user == nil || user.PasswordHash == "" {
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

// FindOrderCustomer returns the existing customer a guest checkout should
// be linked to, or nil. Numbers match in any format. A number that belongs to
// a registered account only links when the checkout email is that account's,
// so registering someone else's number doesn't collect their orders. It never
// creates a record.
func (s *authService) FindOrderCustomer(phone, email string) (*models.User, error) {
	user, err := s.userRepo.FindByPhoneAnyFormat(phone)
	if err != nil || user == nil {
		return nil, err
	}
	if user.PasswordHash != "" && (user.Email == nil || !strings.EqualFold(*user.Email, email)) {
		return nil, nil
	}
	return user, nil
}
