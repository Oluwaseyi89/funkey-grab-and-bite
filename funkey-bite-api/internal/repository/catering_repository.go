package repository

import (
	"database/sql"
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type CateringRepository struct {
	db *sql.DB
}

func NewCateringRepository(db *sql.DB) *CateringRepository {
	return &CateringRepository{db: db}
}

func (r *CateringRepository) Create(request *models.CateringRequest) (*models.CateringRequest, error) {
	query := `
		INSERT INTO catering_requests (
			user_id, event_name, contact_name, contact_phone, contact_email,
			event_date, event_time, guest_count, event_type, package,
			budget, special_requests, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		request.UserID,
		request.EventName,
		request.ContactName,
		request.ContactPhone,
		request.ContactEmail,
		request.EventDate,
		request.EventTime,
		request.GuestCount,
		request.EventType,
		request.Package,
		request.Budget,
		request.SpecialRequests,
		request.Status,
		request.CreatedAt,
	).Scan(&request.ID, &request.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create catering request: %w", err)
	}

	return request, nil
}

func (r *CateringRepository) GetByID(id int) (*models.CateringRequest, error) {
	query := `
		SELECT id, user_id, event_name, contact_name, contact_phone, contact_email,
		       event_date, event_time, guest_count, event_type, package,
		       budget, special_requests, status, created_at
		FROM catering_requests
		WHERE id = $1
	`

	var request models.CateringRequest
	var userID sql.NullInt64
	var eventTime sql.NullString

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&request.ID,
		&userID,
		&request.EventName,
		&request.ContactName,
		&request.ContactPhone,
		&request.ContactEmail,
		&request.EventDate,
		&eventTime,
		&request.GuestCount,
		&request.EventType,
		&request.Package,
		&request.Budget,
		&request.SpecialRequests,
		&request.Status,
		&request.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get catering request: %w", err)
	}

	if userID.Valid {
		val := int(userID.Int64)
		request.UserID = &val
	}
	if eventTime.Valid {
		request.EventTime = &eventTime.String
	}

	return &request, nil
}

func (r *CateringRepository) GetByUserID(userID int) ([]models.CateringRequest, error) {
	query := `
		SELECT id, user_id, event_name, contact_name, contact_phone, contact_email,
		       event_date, event_time, guest_count, event_type, package,
		       budget, special_requests, status, created_at
		FROM catering_requests
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user catering requests: %w", err)
	}
	defer rows.Close()

	var requests []models.CateringRequest
	for rows.Next() {
		var request models.CateringRequest
		var uid sql.NullInt64
		var eventTime sql.NullString

		err := rows.Scan(
			&request.ID,
			&uid,
			&request.EventName,
			&request.ContactName,
			&request.ContactPhone,
			&request.ContactEmail,
			&request.EventDate,
			&eventTime,
			&request.GuestCount,
			&request.EventType,
			&request.Package,
			&request.Budget,
			&request.SpecialRequests,
			&request.Status,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan catering request: %w", err)
		}

		if uid.Valid {
			val := int(uid.Int64)
			request.UserID = &val
		}
		if eventTime.Valid {
			request.EventTime = &eventTime.String
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *CateringRepository) GetAll() ([]models.CateringRequest, error) {
	query := `
		SELECT id, user_id, event_name, contact_name, contact_phone, contact_email,
		       event_date, event_time, guest_count, event_type, package,
		       budget, special_requests, status, created_at
		FROM catering_requests
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all catering requests: %w", err)
	}
	defer rows.Close()

	var requests []models.CateringRequest
	for rows.Next() {
		var request models.CateringRequest
		var userID sql.NullInt64
		var eventTime sql.NullString

		err := rows.Scan(
			&request.ID,
			&userID,
			&request.EventName,
			&request.ContactName,
			&request.ContactPhone,
			&request.ContactEmail,
			&request.EventDate,
			&eventTime,
			&request.GuestCount,
			&request.EventType,
			&request.Package,
			&request.Budget,
			&request.SpecialRequests,
			&request.Status,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan catering request: %w", err)
		}

		if userID.Valid {
			val := int(userID.Int64)
			request.UserID = &val
		}
		if eventTime.Valid {
			request.EventTime = &eventTime.String
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *CateringRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE catering_requests 
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update catering request status: %w", err)
	}

	return nil
}

func (r *CateringRepository) GetByStatus(status string) ([]models.CateringRequest, error) {
	query := `
		SELECT id, user_id, event_name, contact_name, contact_phone, contact_email,
		       event_date, event_time, guest_count, event_type, package,
		       budget, special_requests, status, created_at
		FROM catering_requests
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get catering requests by status: %w", err)
	}
	defer rows.Close()

	var requests []models.CateringRequest
	for rows.Next() {
		var request models.CateringRequest
		var userID sql.NullInt64
		var eventTime sql.NullString

		err := rows.Scan(
			&request.ID,
			&userID,
			&request.EventName,
			&request.ContactName,
			&request.ContactPhone,
			&request.ContactEmail,
			&request.EventDate,
			&eventTime,
			&request.GuestCount,
			&request.EventType,
			&request.Package,
			&request.Budget,
			&request.SpecialRequests,
			&request.Status,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan catering request: %w", err)
		}

		if userID.Valid {
			val := int(userID.Int64)
			request.UserID = &val
		}
		if eventTime.Valid {
			request.EventTime = &eventTime.String
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *CateringRepository) Delete(id int) error {
	query := `DELETE FROM catering_requests WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CateringRepository) GetAllWithPagination(limit, offset int, status string) ([]models.CateringRequest, error) {
	query := `
        SELECT id, user_id, event_name, contact_name, contact_phone, contact_email,
               event_date, event_time, guest_count, event_type, package,
               budget, special_requests, status, created_at
        FROM catering_requests
        WHERE ($1 = '' OR status = $1)
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get catering requests: %w", err)
	}
	defer rows.Close()

	return r.scanCateringRequests(rows)
}

func (r *CateringRepository) GetCount(status string) (int, error) {
	query := `SELECT COUNT(*) FROM catering_requests WHERE ($1 = '' OR status = $1)`

	var count int
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get catering count: %w", err)
	}

	return count, nil
}

func (r *CateringRepository) scanCateringRequests(rows *sql.Rows) ([]models.CateringRequest, error) {
	var requests []models.CateringRequest

	for rows.Next() {
		var request models.CateringRequest
		var userID sql.NullInt64
		var eventTime sql.NullString

		err := rows.Scan(
			&request.ID,
			&userID,
			&request.EventName,
			&request.ContactName,
			&request.ContactPhone,
			&request.ContactEmail,
			&request.EventDate,
			&eventTime,
			&request.GuestCount,
			&request.EventType,
			&request.Package,
			&request.Budget,
			&request.SpecialRequests,
			&request.Status,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan catering request: %w", err)
		}

		if userID.Valid {
			val := int(userID.Int64)
			request.UserID = &val
		}
		if eventTime.Valid {
			request.EventTime = &eventTime.String
		}

		requests = append(requests, request)
	}

	return requests, nil
}
