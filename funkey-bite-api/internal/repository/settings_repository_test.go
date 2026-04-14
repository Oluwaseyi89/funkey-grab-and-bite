package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetSettingsBootstrapsDefaultsWhenNoSettingsExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	repo := NewSettingsRepository(db)
	createdAt := time.Now()
	updatedAt := createdAt.Add(5 * time.Second)

	selectQuery := regexp.QuoteMeta(`
        SELECT id, business_name, phone_number, email, address, opening_hours,
               delivery_fee, min_order_amount, tax_rate, is_delivery_open,
               is_pickup_open, created_at, updated_at
        FROM business_settings
        ORDER BY id DESC
        LIMIT 1
    `)

	mock.ExpectQuery(selectQuery).WillReturnError(sql.ErrNoRows)

	insertQuery := regexp.QuoteMeta(`
        INSERT INTO business_settings (
            business_name, phone_number, email, address, opening_hours,
            delivery_fee, min_order_amount, tax_rate, is_delivery_open,
            is_pickup_open, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at, updated_at
    `)
	mock.ExpectQuery(insertQuery).
		WithArgs(
			"Funkey Grab & Bite",
			"+1234567890",
			"contact@funkeygrabandbite.com",
			"123 Fast Food Street, Food City",
			sqlmock.AnyArg(),
			2.99,
			10.00,
			8.5,
			true,
			true,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, createdAt, updatedAt))

	settings, err := repo.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings() error = %v", err)
	}

	if settings.ID != 1 {
		t.Fatalf("expected created settings ID to be 1, got %d", settings.ID)
	}

	if settings.BusinessName != "Funkey Grab & Bite" {
		t.Fatalf("expected default business name, got %q", settings.BusinessName)
	}

	if len(settings.OpeningHours) != 7 {
		t.Fatalf("expected 7 default opening-hours entries, got %d", len(settings.OpeningHours))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
