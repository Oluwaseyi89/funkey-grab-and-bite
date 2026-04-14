package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestGetPublicSettingsFailsSafelyWhenOpeningHoursJSONIsMalformed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	selectQuery := regexp.QuoteMeta(`
        SELECT id, business_name, phone_number, email, address, opening_hours,
               delivery_fee, min_order_amount, tax_rate, is_delivery_open,
               is_pickup_open, created_at, updated_at
        FROM business_settings
        ORDER BY id DESC
        LIMIT 1
    `)

	mock.ExpectQuery(selectQuery).WillReturnRows(
		sqlmock.NewRows([]string{
			"id", "business_name", "phone_number", "email", "address", "opening_hours",
			"delivery_fee", "min_order_amount", "tax_rate", "is_delivery_open",
			"is_pickup_open", "created_at", "updated_at",
		}).AddRow(
			1,
			"Funkey Grab & Bite",
			"+1234567890",
			"contact@funkeygrabandbite.com",
			"123 Fast Food Street, Food City",
			"{not-valid-json",
			2.99,
			10.00,
			8.5,
			true,
			true,
			time.Now(),
			time.Now(),
		),
	)

	settingsRepo := repository.NewSettingsRepository(db)
	settingsSvc := services.NewSettingsService(*settingsRepo)
	handler := NewSettingsHandler(settingsSvc)

	router := gin.New()
	router.GET("/settings", handler.GetPublicSettings)

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d body=%s", http.StatusInternalServerError, res.Code, res.Body.String())
	}

	var payload handlers.APIResponse
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v body=%s", err, res.Body.String())
	}

	if payload.Success {
		t.Fatalf("expected success=false, got true body=%s", res.Body.String())
	}

	if payload.Error == nil {
		t.Fatalf("expected error payload, got nil body=%s", res.Body.String())
	}

	if payload.Error.Code != "SETTINGS_FETCH_FAILED" {
		t.Fatalf("expected error code SETTINGS_FETCH_FAILED, got %q", payload.Error.Code)
	}

	if payload.Error.Message != "Failed to fetch settings" {
		t.Fatalf("expected generic safe error message, got %q", payload.Error.Message)
	}

	if payload.Error.Details == "" || !strings.Contains(payload.Error.Details, "failed to parse opening hours JSON") {
		t.Fatalf("expected parse failure details, got %q", payload.Error.Details)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
