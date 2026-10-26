package services

import (
	"errors"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

var userColumns = []string{
	"id", "phone", "email", "full_name", "password_hash", "is_verified",
	"is_active", "last_login", "created_at", "updated_at",
}

func newGuestAuthService(t *testing.T) (AuthService, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return NewAuthService(*repository.NewUserRepository(db)), mock
}

func TestResolveOrderCustomerReusesExistingPhone(t *testing.T) {
	svc, mock := newGuestAuthService(t)
	now := time.Now()
	mock.ExpectQuery("FROM users").WithArgs("08012345678").WillReturnRows(
		sqlmock.NewRows(userColumns).AddRow(4, "08012345678", nil, "Ada", "", false, true, nil, now, now),
	)

	user, err := svc.ResolveOrderCustomer("Ada Lovelace", "08012345678", nil)
	if err != nil || user == nil || user.ID != 4 {
		t.Fatalf("expected existing customer 4, got %+v, %v", user, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestResolveOrderCustomerCreatesPasswordlessGuest(t *testing.T) {
	svc, mock := newGuestAuthService(t)
	email := "ada@example.com"
	mock.ExpectQuery("FROM users").WithArgs("08012345678").WillReturnRows(sqlmock.NewRows(userColumns))
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("08012345678", &email, "Ada", "", false, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(11, time.Now(), time.Now()))

	user, err := svc.ResolveOrderCustomer("Ada", "08012345678", &email)
	if err != nil || user.ID != 11 || user.PasswordHash != "" {
		t.Fatalf("expected new passwordless guest 11, got %+v, %v", user, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestResolveOrderCustomerRetriesWithoutTakenEmail(t *testing.T) {
	svc, mock := newGuestAuthService(t)
	email := "shared@example.com"
	mock.ExpectQuery("FROM users").WillReturnRows(sqlmock.NewRows(userColumns))
	mock.ExpectQuery("INSERT INTO users").WillReturnError(errors.New("duplicate key value violates unique constraint \"users_email_key\""))
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("08012345678", nil, "Ada", "", false, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(12, time.Now(), time.Now()))

	user, err := svc.ResolveOrderCustomer("Ada", "08012345678", &email)
	if err != nil || user.ID != 12 || user.Email != nil {
		t.Fatalf("expected guest 12 without email, got %+v, %v", user, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterClaimsGuestAccount(t *testing.T) {
	if err := utils.ConfigureJWTSecret("test-secret-for-guest-claim"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { utils.SetSecret("") })
	svc, mock := newGuestAuthService(t)
	now := time.Now()
	mock.ExpectQuery("FROM users").WillReturnRows(
		sqlmock.NewRows(userColumns).AddRow(4, "08012345678", nil, "Ada", "", false, true, nil, now, now),
	)
	mock.ExpectExec("UPDATE users SET password_hash").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE users SET updated_at").WillReturnResult(sqlmock.NewResult(0, 1))

	resp, err := svc.Register(models.UserRegistration{
		Phone: "08012345678", Email: "ada@example.com", FullName: "Ada Lovelace", Password: "a-strong-password",
	})
	if err != nil || resp.User.ID != 4 || resp.AccessToken == "" {
		t.Fatalf("expected guest 4 to be claimed with a token, got %+v, %v", resp, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestLoginRejectsUnknownPhoneWithoutPanicking(t *testing.T) {
	svc, mock := newGuestAuthService(t)
	mock.ExpectQuery("FROM users").WillReturnRows(sqlmock.NewRows(userColumns))

	if _, err := svc.Login("08000000000", "whatever"); err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
