package database

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRunMigrationsWithStatementsFailsOnConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	migrations := []string{
		"CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY)",
		"ALTER TABLE users ADD COLUMN email TEXT NOT NULL",
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ALTER TABLE users ADD COLUMN email TEXT NOT NULL").WillReturnError(fmt.Errorf("column \"email\" of relation \"users\" already exists"))

	err = runMigrationsWithStatements(db, migrations)
	if err == nil {
		t.Fatal("runMigrationsWithStatements() expected error, got nil")
	}

	if !strings.Contains(err.Error(), "migration 2 failed") {
		t.Fatalf("expected migration index in error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEnsureDefaultAdminUserReturnsErrorWhenCountLookupFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(1\) FROM admin_users`).WillReturnError(fmt.Errorf("relation \"admin_users\" does not exist"))

	err = ensureDefaultAdminUser(db)
	if err == nil {
		t.Fatal("ensureDefaultAdminUser() expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed checking admin user count") {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEnsureDefaultAdminUserFailsWhenDefaultsAreUnset(t *testing.T) {
	t.Setenv("DEFAULT_ADMIN_EMAIL", "")
	t.Setenv("DEFAULT_ADMIN_USERNAME", "")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "")
	t.Setenv("ENVIRONMENT", "development")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(1\) FROM admin_users`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err = ensureDefaultAdminUser(db)
	if err == nil {
		t.Fatal("ensureDefaultAdminUser() expected error when defaults are unset, got nil")
	}

	if !strings.Contains(err.Error(), "missing required default admin bootstrap env vars") {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEnsureDefaultAdminUserRejectsWeakDefaultsInProductionLikeEnv(t *testing.T) {
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@funkey.com")
	t.Setenv("DEFAULT_ADMIN_USERNAME", "admin")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "admin123")
	t.Setenv("ENVIRONMENT", "production")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(1\) FROM admin_users`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err = ensureDefaultAdminUser(db)
	if err == nil {
		t.Fatal("ensureDefaultAdminUser() expected error for weak production defaults, got nil")
	}

	if !strings.Contains(err.Error(), "weak default admin credentials are not allowed") {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEnsureDefaultAdminUserSkipsBootstrapWhenAdminAlreadyExists(t *testing.T) {
	t.Setenv("DEFAULT_ADMIN_EMAIL", "secure-admin@funkey.com")
	t.Setenv("DEFAULT_ADMIN_USERNAME", "secureadmin")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "super-secure-password-123")
	t.Setenv("ENVIRONMENT", "production")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(1\) FROM admin_users`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err = ensureDefaultAdminUser(db)
	if err != nil {
		t.Fatalf("ensureDefaultAdminUser() expected nil when admin exists, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
