package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type fakeAccountUsers struct {
	byPhone     map[string]*models.User
	byEmail     map[string]*models.User
	orderEmail  string
	created     *models.User
	claimed     bool
	newPassword string
}

func (f *fakeAccountUsers) FindByPhoneAnyFormat(phone string) (*models.User, error) {
	return f.byPhone[utils.CanonicalPhone(phone)], nil
}
func (f *fakeAccountUsers) FindByEmail(email string) (*models.User, error) {
	return f.byEmail[strings.ToLower(email)], nil
}
func (f *fakeAccountUsers) LatestOrderEmail(userID int) (string, error) { return f.orderEmail, nil }
func (f *fakeAccountUsers) Create(user *models.User) (*models.User, error) {
	user.ID = 50
	f.created = user
	return user, nil
}
func (f *fakeAccountUsers) ClaimGuest(userID int, fullName, email, passwordHash string) (bool, error) {
	f.claimed = true
	for _, u := range f.byPhone {
		if u.ID == userID {
			u.PasswordHash, u.FullName, u.Email = passwordHash, fullName, &email
		}
	}
	return true, nil
}
func (f *fakeAccountUsers) SetPassword(userID int, passwordHash string) error {
	f.newPassword = passwordHash
	return nil
}

type fakeAccountCodes struct {
	hashes map[string]string
}

func (f *fakeAccountCodes) Create(ctx context.Context, userID int, purpose, email, codeHash string, ttl time.Duration) error {
	if f.hashes == nil {
		f.hashes = map[string]string{}
	}
	f.hashes[purpose] = codeHash
	return nil
}
func (f *fakeAccountCodes) Use(ctx context.Context, userID int, purpose, codeHash string, maxAttempts int) (bool, error) {
	if f.hashes[purpose] != "" && f.hashes[purpose] == codeHash {
		delete(f.hashes, purpose)
		return true, nil
	}
	return false, nil
}

type fakeEmailQueue struct{ sent []models.EmailJobPayload }

func (f *fakeEmailQueue) QueueEmail(ctx context.Context, email models.EmailJobPayload) error {
	f.sent = append(f.sent, email)
	return nil
}

type fakeAccountAudit struct{ actions []string }

func (f *fakeAccountAudit) Record(ctx context.Context, exec repository.Execer, entry models.AuditEntry) error {
	f.actions = append(f.actions, entry.Action)
	return nil
}

func newTestAccounts(t *testing.T, users *fakeAccountUsers) (*accountService, *fakeEmailQueue, *fakeAccountAudit) {
	t.Helper()
	if err := utils.ConfigureJWTSecret("test-secret-for-accounts-0123456789"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { utils.SetSecret("") })
	queue, audit := &fakeEmailQueue{}, &fakeAccountAudit{}
	svc := NewAccountService(users, &fakeAccountCodes{}, queue, audit).(*accountService)
	svc.newCode = func() (string, error) { return "123456", nil }
	return svc, queue, audit
}

func registration(phone, email string) models.UserRegistration {
	return models.UserRegistration{FullName: "Ada Lovelace", Phone: phone, Email: email, Password: "Secret123"}
}

func strptr(s string) *string { return &s }

func TestRegisterCreatesAccountForNewNumber(t *testing.T) {
	users := &fakeAccountUsers{}
	svc, queue, _ := newTestAccounts(t, users)

	result, err := svc.Register(context.Background(), registration("0803 111 2222", "ada@example.test"))
	if err != nil || result.Auth == nil || result.Auth.AccessToken == "" || !result.Auth.IsNewUser {
		t.Fatalf("expected a signed-in new account, got %+v %v", result, err)
	}
	if users.created.PasswordHash == "" || users.created.IsVerified {
		t.Fatal("a new account has a password and isn't marked verified")
	}
	if len(queue.sent) != 0 {
		t.Fatal("no code is needed for a number without guest orders")
	}
}

func TestRegisterOverGuestOrdersNeedsTheEmailedCode(t *testing.T) {
	guest := &models.User{ID: 7, Phone: "+2348031112222", Email: strptr("ada@example.test"), FullName: "Ada"}
	users := &fakeAccountUsers{byPhone: map[string]*models.User{"08031112222": guest}, byEmail: map[string]*models.User{"ada@example.test": guest}}
	svc, queue, audit := newTestAccounts(t, users)
	ctx := context.Background()

	result, err := svc.Register(ctx, registration("0803 111 2222", "ada@example.test"))
	if err != nil || result.Auth != nil || result.VerificationSentTo != "a**@example.test" {
		t.Fatalf("expected a code sent instead of an account, got %+v %v", result, err)
	}
	if len(queue.sent) != 1 || queue.sent[0].To != "ada@example.test" || !strings.Contains(queue.sent[0].Text, "123456") {
		t.Fatalf("expected the code emailed to the guest's address, got %+v", queue.sent)
	}
	if users.claimed {
		t.Fatal("the guest record must not be taken over before the code is confirmed")
	}

	if _, err := svc.ConfirmRegistration(ctx, registration("0803 111 2222", "ada@example.test"), "000000"); !errors.Is(err, ErrInvalidCode) {
		t.Fatalf("a wrong code must be refused, got %v", err)
	}
	auth, err := svc.ConfirmRegistration(ctx, registration("+234 803 111 2222", "ada@example.test"), "123456")
	if err != nil || auth.AccessToken == "" || !users.claimed {
		t.Fatalf("expected the record claimed and signed in, got %+v %v", auth, err)
	}
	if len(audit.actions) != 1 || audit.actions[0] != "customer.account_claimed" {
		t.Fatalf("expected the claim audited, got %v", audit.actions)
	}
}

func TestRegisterRefusals(t *testing.T) {
	registered := &models.User{ID: 3, Phone: "08035556666", PasswordHash: "hash", Email: strptr("taken@example.test")}
	noEmailGuest := &models.User{ID: 9, Phone: "08037778888"}
	users := &fakeAccountUsers{
		byPhone: map[string]*models.User{"08035556666": registered, "08037778888": noEmailGuest},
		byEmail: map[string]*models.User{"taken@example.test": registered},
	}
	svc, _, _ := newTestAccounts(t, users)
	ctx := context.Background()

	if _, err := svc.Register(ctx, registration("08035556666", "new@example.test")); !errors.Is(err, ErrAccountExists) {
		t.Fatalf("a registered number must be refused, got %v", err)
	}
	if _, err := svc.Register(ctx, registration("08039990000", "taken@example.test")); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("another account's email must be refused, got %v", err)
	}
	if _, err := svc.Register(ctx, registration("08037778888", "me@example.test")); !errors.Is(err, ErrCannotVerifyPhone) {
		t.Fatalf("a guest record with no email to verify must be refused, got %v", err)
	}
}

func TestRegisterUsesLatestOrderEmailForGuestWithoutOne(t *testing.T) {
	guest := &models.User{ID: 9, Phone: "08037778888"}
	users := &fakeAccountUsers{byPhone: map[string]*models.User{"08037778888": guest}, orderEmail: "orders@example.test"}
	svc, queue, _ := newTestAccounts(t, users)

	if _, err := svc.Register(context.Background(), registration("08037778888", "me@example.test")); err != nil {
		t.Fatal(err)
	}
	if len(queue.sent) != 1 || queue.sent[0].To != "orders@example.test" {
		t.Fatalf("expected the code sent to the email on the latest order, got %+v", queue.sent)
	}
}

func TestPasswordReset(t *testing.T) {
	user := &models.User{ID: 3, Phone: "08035556666", PasswordHash: "old", Email: strptr("ada@example.test"), IsActive: true}
	users := &fakeAccountUsers{byEmail: map[string]*models.User{"ada@example.test": user}}
	svc, queue, audit := newTestAccounts(t, users)
	ctx := context.Background()

	if err := svc.RequestPasswordReset(ctx, "nobody@example.test"); err != nil || len(queue.sent) != 0 {
		t.Fatalf("an unknown email gets no code and no error, got %v and %d emails", err, len(queue.sent))
	}
	if err := svc.RequestPasswordReset(ctx, "ADA@example.test"); err != nil || len(queue.sent) != 1 {
		t.Fatalf("expected a reset code emailed, got %v and %d emails", err, len(queue.sent))
	}
	if err := svc.ResetPassword(ctx, "ada@example.test", "111111", "NewSecret1"); !errors.Is(err, ErrInvalidCode) {
		t.Fatalf("a wrong code must be refused, got %v", err)
	}
	if err := svc.ResetPassword(ctx, "ada@example.test", "123456", "NewSecret1"); err != nil {
		t.Fatal(err)
	}
	if !utils.CheckPasswordHash("NewSecret1", users.newPassword) {
		t.Fatal("expected the new password saved")
	}
	if err := svc.ResetPassword(ctx, "ada@example.test", "123456", "Another1A"); !errors.Is(err, ErrInvalidCode) {
		t.Fatalf("a code works only once, got %v", err)
	}
	if len(audit.actions) != 1 || audit.actions[0] != "customer.password_reset" {
		t.Fatalf("expected the reset audited, got %v", audit.actions)
	}
}

func TestCodesAreKeyedToUserAndPurpose(t *testing.T) {
	if err := utils.ConfigureJWTSecret("test-secret-for-accounts-0123456789"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { utils.SetSecret("") })
	a := hashCode(1, repository.CodePurposeClaimAccount, "123456")
	if a == hashCode(2, repository.CodePurposeClaimAccount, "123456") || a == hashCode(1, repository.CodePurposePasswordReset, "123456") {
		t.Fatal("the same code must hash differently for another user or purpose")
	}
}
