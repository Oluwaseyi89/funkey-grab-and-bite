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

func TestEnsureDefaultAdminUserReturnsErrorWhenLookupFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM admin_users WHERE email = \$1\)`).WithArgs("admin@funkey.com").WillReturnError(fmt.Errorf("relation \"admin_users\" does not exist"))

	err = ensureDefaultAdminUser(db)
	if err == nil {
		t.Fatal("ensureDefaultAdminUser() expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed checking default admin existence") {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
