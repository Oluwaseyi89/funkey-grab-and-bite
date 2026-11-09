package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

var (
	ErrAccountExists     = errors.New("an account with this phone number already exists; log in or reset your password")
	ErrEmailTaken        = errors.New("this email is already used by another account; log in or use a different email")
	ErrCannotVerifyPhone = errors.New("we can't confirm this phone number online; please contact us to set up your account")
	ErrInvalidCode       = errors.New("that code is wrong or has expired; ask for a new one")
)

const (
	codeTTL         = 15 * time.Minute
	codeMaxAttempts = 5
)

// RegisterResult is either a signed-in account, or (when the number already
// has guest orders) where the code to confirm it was sent.
type RegisterResult struct {
	Auth               *models.AuthResponse
	VerificationSentTo string
}

// AccountService registers customers and resets their passwords.
type AccountService interface {
	// Register creates an account. A number with guest orders isn't taken
	// over yet: a code is emailed to the address on those orders, to be
	// confirmed with ConfirmRegistration.
	Register(ctx context.Context, reg models.UserRegistration) (*RegisterResult, error)
	// ConfirmRegistration takes over the guest record, with its order
	// history, once the emailed code is entered.
	ConfirmRegistration(ctx context.Context, reg models.UserRegistration, code string) (*models.AuthResponse, error)
	// RequestPasswordReset emails a reset code to an account with this email.
	// It says nothing about whether one exists.
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, code, password string) error
}

type accountUsers interface {
	FindByPhoneAnyFormat(phone string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	LatestOrderEmail(userID int) (string, error)
	Create(user *models.User) (*models.User, error)
	ClaimGuest(userID int, fullName, email, passwordHash string) (bool, error)
	SetPassword(userID int, passwordHash string) error
}

type accountCodes interface {
	Create(ctx context.Context, userID int, purpose, email, codeHash string, ttl time.Duration) error
	Use(ctx context.Context, userID int, purpose, codeHash string, maxAttempts int) (bool, error)
}

type accountAudit interface {
	Record(ctx context.Context, exec repository.Execer, entry models.AuditEntry) error
}

type accountService struct {
	users accountUsers
	codes accountCodes
	email EmailQueue
	audit accountAudit
	// newCode makes a code; tests replace it.
	newCode func() (string, error)
}

func NewAccountService(users accountUsers, codes accountCodes, email EmailQueue, audit accountAudit) AccountService {
	return &accountService{users: users, codes: codes, email: email, audit: audit, newCode: sixDigitCode}
}

func sixDigitCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// hashCode keys the code to its user and purpose with the server secret, so
// a code is useless for anything else and the stored hash can't be reversed
// by trying all million codes without the secret.
func hashCode(userID int, purpose, code string) string {
	mac := hmac.New(sha256.New, utils.GetSecret())
	mac.Write([]byte(purpose + ":" + strconv.Itoa(userID) + ":" + strings.TrimSpace(code)))
	return hex.EncodeToString(mac.Sum(nil))
}

// maskEmail shows enough of an address for its owner to recognise it.
func maskEmail(email string) string {
	local, domain, ok := strings.Cut(email, "@")
	if !ok || local == "" {
		return "your email"
	}
	return local[:1] + strings.Repeat("*", max(len(local)-1, 2)) + "@" + domain
}

func (s *accountService) sendCode(ctx context.Context, user *models.User, purpose, to string) error {
	code, err := s.newCode()
	if err != nil {
		return fmt.Errorf("make code: %w", err)
	}
	if err := s.codes.Create(ctx, user.ID, purpose, to, hashCode(user.ID, purpose, code), codeTTL); err != nil {
		return err
	}

	what := "confirm your account"
	if purpose == repository.CodePurposePasswordReset {
		what = "reset your password"
	}
	text := fmt.Sprintf("Your code to %s is %s. It expires in 15 minutes.\n\nIf you didn't ask for this, you can ignore this email.", what, code)
	html := fmt.Sprintf(`<p>Your code to %s is</p><p style="font-size:24px;font-weight:bold;letter-spacing:4px">%s</p><p>It expires in 15 minutes.</p><p>If you didn't ask for this, you can ignore this email.</p>`, what, code)
	return s.email.QueueEmail(ctx, models.EmailJobPayload{To: to, Subject: "Your verification code", HTML: html, Text: text})
}

// emailTakenByOther reports whether email belongs to a customer other than
// userID (0 for a new customer).
func (s *accountService) emailTakenByOther(email string, userID int) (bool, error) {
	other, err := s.users.FindByEmail(email)
	if err != nil {
		return false, err
	}
	return other != nil && other.ID != userID, nil
}

func (s *accountService) Register(ctx context.Context, reg models.UserRegistration) (*RegisterResult, error) {
	phone := utils.CleanPhone(reg.Phone)
	email := strings.TrimSpace(reg.Email)

	existing, err := s.users.FindByPhoneAnyFormat(phone)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.PasswordHash != "" {
		return nil, ErrAccountExists
	}
	existingID := 0
	if existing != nil {
		existingID = existing.ID
	}
	if taken, err := s.emailTakenByOther(email, existingID); err != nil {
		return nil, err
	} else if taken {
		return nil, ErrEmailTaken
	}

	if existing != nil {
		// The number has guest orders: prove it's yours through the email
		// those orders were placed with before seeing them.
		to := ""
		if existing.Email != nil {
			to = *existing.Email
		}
		if to == "" {
			if to, err = s.users.LatestOrderEmail(existing.ID); err != nil {
				return nil, err
			}
		}
		if to == "" {
			return nil, ErrCannotVerifyPhone
		}
		if err := s.sendCode(ctx, existing, repository.CodePurposeClaimAccount, to); err != nil {
			return nil, err
		}
		return &RegisterResult{VerificationSentTo: maskEmail(to)}, nil
	}

	hash, err := utils.HashPassword(reg.Password)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user, err := s.users.Create(&models.User{
		Phone: phone, Email: &email, FullName: strings.TrimSpace(reg.FullName), PasswordHash: hash,
		IsActive: true, CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	auth, err := s.signedIn(user, true)
	if err != nil {
		return nil, err
	}
	return &RegisterResult{Auth: auth}, nil
}

func (s *accountService) ConfirmRegistration(ctx context.Context, reg models.UserRegistration, code string) (*models.AuthResponse, error) {
	phone := utils.CleanPhone(reg.Phone)
	email := strings.TrimSpace(reg.Email)

	existing, err := s.users.FindByPhoneAnyFormat(phone)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrInvalidCode
	}
	if existing.PasswordHash != "" {
		return nil, ErrAccountExists
	}
	ok, err := s.codes.Use(ctx, existing.ID, repository.CodePurposeClaimAccount, hashCode(existing.ID, repository.CodePurposeClaimAccount, code), codeMaxAttempts)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCode
	}
	if taken, err := s.emailTakenByOther(email, existing.ID); err != nil {
		return nil, err
	} else if taken {
		return nil, ErrEmailTaken
	}

	hash, err := utils.HashPassword(reg.Password)
	if err != nil {
		return nil, err
	}
	claimed, err := s.users.ClaimGuest(existing.ID, strings.TrimSpace(reg.FullName), email, hash)
	if err != nil {
		return nil, err
	}
	if !claimed {
		return nil, ErrAccountExists
	}

	userCtx := withCustomer(ctx, existing.ID)
	if err := s.audit.Record(userCtx, nil, models.AuditEntry{
		Action: "customer.account_claimed", EntityType: "customer", EntityID: strconv.Itoa(existing.ID),
		Metadata: map[string]any{"verifiedBy": "email code"},
	}); err != nil {
		return nil, err
	}

	user, err := s.users.FindByPhoneAnyFormat(phone)
	if err != nil {
		return nil, err
	}
	return s.signedIn(user, false)
}

func (s *accountService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.users.FindByEmail(strings.TrimSpace(email))
	if err != nil {
		return err
	}
	if user == nil || user.PasswordHash == "" || !user.IsActive || user.Email == nil {
		return nil
	}
	return s.sendCode(ctx, user, repository.CodePurposePasswordReset, *user.Email)
}

func (s *accountService) ResetPassword(ctx context.Context, email, code, password string) error {
	user, err := s.users.FindByEmail(strings.TrimSpace(email))
	if err != nil {
		return err
	}
	if user == nil || user.PasswordHash == "" {
		return ErrInvalidCode
	}
	ok, err := s.codes.Use(ctx, user.ID, repository.CodePurposePasswordReset, hashCode(user.ID, repository.CodePurposePasswordReset, code), codeMaxAttempts)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidCode
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	if err := s.users.SetPassword(user.ID, hash); err != nil {
		return err
	}
	return s.audit.Record(withCustomer(ctx, user.ID), nil, models.AuditEntry{
		Action: "customer.password_reset", EntityType: "customer", EntityID: strconv.Itoa(user.ID),
	})
}

func (s *accountService) signedIn(user *models.User, isNew bool) (*models.AuthResponse, error) {
	token, err := utils.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{User: user, AccessToken: token, IsNewUser: isNew}, nil
}

// withCustomer records the customer as the actor of what follows, keeping
// the request's IP and ID.
func withCustomer(ctx context.Context, userID int) context.Context {
	actor := requestctx.ActorFrom(ctx)
	actor.Type, actor.ID = requestctx.ActorCustomer, strconv.Itoa(userID)
	return requestctx.WithActor(ctx, actor)
}
