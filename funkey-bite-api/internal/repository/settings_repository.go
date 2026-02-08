package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// GetSettings gets the business settings
func (r *SettingsRepository) GetSettings() (*models.BusinessSettings, error) {
	query := `
        SELECT id, business_name, phone_number, email, address, opening_hours,
               delivery_fee, min_order_amount, tax_rate, is_delivery_open,
               is_pickup_open, created_at, updated_at
        FROM business_settings
        ORDER BY id DESC
        LIMIT 1
    `

	var settings models.BusinessSettings
	var openingHoursJSON string // Store JSON as string from DB

	row := r.db.QueryRow(query)

	err := row.Scan(
		&settings.ID,
		&settings.BusinessName,
		&settings.PhoneNumber,
		&settings.Email,
		&settings.Address,
		&openingHoursJSON, // Scan into string
		&settings.DeliveryFee,
		&settings.MinOrderAmount,
		&settings.TaxRate,
		&settings.IsDeliveryOpen,
		&settings.IsPickupOpen,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default settings if none exist
		return r.createDefaultSettings()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// Parse JSON into slice
	if err := json.Unmarshal([]byte(openingHoursJSON), &settings.OpeningHours); err != nil {
		return nil, fmt.Errorf("failed to parse opening hours JSON: %w", err)
	}

	return &settings, nil
}

// UpdateSettings updates business settings
func (r *SettingsRepository) UpdateSettings(updates *models.BusinessSettingsUpdate) (*models.BusinessSettings, error) {
	// Start by getting current settings
	current, err := r.GetSettings()
	if err != nil {
		return nil, err
	}

	// Build update query dynamically
	query := "UPDATE business_settings SET "
	params := []interface{}{}
	paramCount := 1

	if updates.BusinessName != nil {
		query += fmt.Sprintf("business_name = $%d, ", paramCount)
		params = append(params, *updates.BusinessName)
		paramCount++
	}
	if updates.PhoneNumber != nil {
		query += fmt.Sprintf("phone_number = $%d, ", paramCount)
		params = append(params, *updates.PhoneNumber)
		paramCount++
	}
	if updates.Email != nil {
		query += fmt.Sprintf("email = $%d, ", paramCount)
		params = append(params, *updates.Email)
		paramCount++
	}
	if updates.Address != nil {
		query += fmt.Sprintf("address = $%d, ", paramCount)
		params = append(params, *updates.Address)
		paramCount++
	}
	if updates.OpeningHours != nil {
		// Convert slice to JSON string
		openingHoursJSON, err := json.Marshal(updates.OpeningHours)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal opening hours: %w", err)
		}
		query += fmt.Sprintf("opening_hours = $%d, ", paramCount)
		params = append(params, string(openingHoursJSON))
		paramCount++
	}
	if updates.DeliveryFee != nil {
		query += fmt.Sprintf("delivery_fee = $%d, ", paramCount)
		params = append(params, *updates.DeliveryFee)
		paramCount++
	}
	if updates.MinOrderAmount != nil {
		query += fmt.Sprintf("min_order_amount = $%d, ", paramCount)
		params = append(params, *updates.MinOrderAmount)
		paramCount++
	}
	if updates.TaxRate != nil {
		query += fmt.Sprintf("tax_rate = $%d, ", paramCount)
		params = append(params, *updates.TaxRate)
		paramCount++
	}
	if updates.IsDeliveryOpen != nil {
		query += fmt.Sprintf("is_delivery_open = $%d, ", paramCount)
		params = append(params, *updates.IsDeliveryOpen)
		paramCount++
	}
	if updates.IsPickupOpen != nil {
		query += fmt.Sprintf("is_pickup_open = $%d, ", paramCount)
		params = append(params, *updates.IsPickupOpen)
		paramCount++
	}

	// Add updated_at and WHERE clause
	query += fmt.Sprintf("updated_at = $%d WHERE id = $%d", paramCount, paramCount+1)
	params = append(params, time.Now(), current.ID)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	// Return updated settings
	return r.GetSettings()
}

// createDefaultSettings creates default business settings
func (r *SettingsRepository) createDefaultSettings() (*models.BusinessSettings, error) {
	defaultHours := []models.OpeningHours{
		{Day: "Monday", OpenTime: "08:00", CloseTime: "22:00", IsOpen: true},
		{Day: "Tuesday", OpenTime: "08:00", CloseTime: "22:00", IsOpen: true},
		{Day: "Wednesday", OpenTime: "08:00", CloseTime: "22:00", IsOpen: true},
		{Day: "Thursday", OpenTime: "08:00", CloseTime: "22:00", IsOpen: true},
		{Day: "Friday", OpenTime: "08:00", CloseTime: "23:00", IsOpen: true},
		{Day: "Saturday", OpenTime: "09:00", CloseTime: "23:00", IsOpen: true},
		{Day: "Sunday", OpenTime: "10:00", CloseTime: "21:00", IsOpen: true},
	}

	// Convert to JSON for database storage
	openingHoursJSON, err := json.Marshal(defaultHours)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default opening hours: %w", err)
	}

	query := `
        INSERT INTO business_settings (
            business_name, phone_number, email, address, opening_hours,
            delivery_fee, min_order_amount, tax_rate, is_delivery_open,
            is_pickup_open, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at, updated_at
    `

	now := time.Now()
	var settings models.BusinessSettings

	err = r.db.QueryRow(
		query,
		"Funkey Grab & Bite",
		"+1234567890",
		"contact@funkeygrabandbite.com",
		"123 Fast Food Street, Food City",
		string(openingHoursJSON),
		2.99,  // delivery_fee
		10.00, // min_order_amount
		8.5,   // tax_rate
		true,  // is_delivery_open
		true,  // is_pickup_open
		now,
		now,
	).Scan(&settings.ID, &settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create default settings: %w", err)
	}

	// Set the other fields
	settings.BusinessName = "Funkey Grab & Bite"
	settings.PhoneNumber = "+1234567890"
	settings.Email = "contact@funkeygrabandbite.com"
	settings.Address = "123 Fast Food Street, Food City"
	settings.OpeningHours = defaultHours // Use the slice, not JSON string
	settings.DeliveryFee = 2.99
	settings.MinOrderAmount = 10.00
	settings.TaxRate = 8.5
	settings.IsDeliveryOpen = true
	settings.IsPickupOpen = true

	return &settings, nil
}
