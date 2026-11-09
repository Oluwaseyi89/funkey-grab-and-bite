package services

import (
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
