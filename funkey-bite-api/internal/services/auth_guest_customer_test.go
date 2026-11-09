package services

import (
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/repository"

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

func TestLoginRejectsUnknownPhoneWithoutPanicking(t *testing.T) {
	svc, mock := newGuestAuthService(t)
	mock.ExpectQuery("FROM users").WillReturnRows(sqlmock.NewRows(userColumns))

	if _, err := svc.Login("08000000000", "whatever"); err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
