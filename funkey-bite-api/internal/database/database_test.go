package database

import (
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/migrations"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEmbeddedMigrationsArePairedAndContiguous(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		t.Fatalf("reading embedded migrations: %v", err)
	}

	filePattern := regexp.MustCompile(`^(\d{6})_[a-z0-9_]+\.(up|down)\.sql$`)
	directions := map[int]map[string]bool{}
	for _, entry := range entries {
		match := filePattern.FindStringSubmatch(entry.Name())
		if match == nil {
			t.Fatalf("migration %q does not match NNNNNN_description.(up|down).sql", entry.Name())
		}
		version, _ := strconv.Atoi(match[1])
		if directions[version] == nil {
			directions[version] = map[string]bool{}
		}
		directions[version][match[2]] = true
	}

	if len(directions) == 0 {
		t.Fatal("expected embedded migrations, found none")
	}

	for version := 1; version <= len(directions); version++ {
		dirs, ok := directions[version]
		if !ok {
			t.Fatalf("migration versions are not contiguous: missing %06d", version)
		}
		if !dirs["up"] || !dirs["down"] {
			t.Fatalf("migration %06d must have both up and down files, got %v", version, dirs)
		}
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
