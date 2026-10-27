package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/lib/pq"
)

// ErrDuplicateCateringReference is returned when a generated reference is
// already taken, so the caller can retry with a new one.
var ErrDuplicateCateringReference = errors.New("catering reference already exists")

type CateringRepository struct {
	db *sql.DB
}

func NewCateringRepository(db *sql.DB) *CateringRepository {
	return &CateringRepository{db: db}
}

// cateringSelect reads full requests; packageName falls back to the free-text
// package stored on requests made before packages were managed.
const cateringSelect = `
	SELECT cr.id, cr.reference, cr.request_type, cr.user_id, cr.package_id,
	       COALESCE(cp.name, cr.package),
	       cr.contact_name, cr.contact_phone, cr.contact_email, cr.organization,
	       cr.event_type, cr.event_type_other, cr.event_name, cr.event_date,
	       cr.event_time, cr.duration_hours, cr.guest_count,
	       cr.venue_name, cr.venue_address, cr.venue_area, cr.venue_state, cr.venue_notes,
	       cr.service_style, cr.selected_items, cr.menu_notes,
	       cr.dietary_requirements, cr.dietary_notes, cr.budget, cr.special_requests,
	       cr.estimated_total, cr.quoted_amount, cr.admin_notes, cr.status,
	       cr.created_at, cr.updated_at, cr.accepted_quote_id, cr.invoice_number,
	       cr.contracted_at, cr.paid_in_full_at
	FROM catering_requests cr
	LEFT JOIN catering_packages cp ON cp.id = cr.package_id
`

func (r *CateringRepository) Create(request *models.CateringRequest) (*models.CateringRequest, error) {
	selectedItems, err := json.Marshal(request.SelectedItems)
	if err != nil {
		return nil, fmt.Errorf("failed to encode selected items: %w", err)
	}
	dietary, err := json.Marshal(request.DietaryRequirements)
	if err != nil {
		return nil, fmt.Errorf("failed to encode dietary requirements: %w", err)
	}

	query := `
		INSERT INTO catering_requests (
			reference, request_type, user_id, package_id, contact_name, contact_phone,
			contact_email, organization, event_type, event_type_other, event_name,
			event_date, event_time, duration_hours, guest_count, venue_name,
			venue_address, venue_area, venue_state, venue_notes, service_style,
			selected_items, menu_notes, dietary_requirements, dietary_notes, budget,
			special_requests, estimated_total, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		        $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $30)
		RETURNING id
	`

	err = r.db.QueryRow(
		query,
		request.Reference, request.RequestType, request.UserID, request.PackageID,
		request.ContactName, request.ContactPhone, request.ContactEmail, request.Organization,
		request.EventType, request.EventTypeOther, request.EventName, request.EventDate,
		request.EventTime, request.DurationHours, request.GuestCount, request.VenueName,
		request.VenueAddress, request.VenueArea, request.VenueState, request.VenueNotes,
		request.ServiceStyle, selectedItems, request.MenuNotes, dietary, request.DietaryNotes,
		request.Budget, request.SpecialRequests, request.EstimatedTotal, request.Status,
		request.CreatedAt,
	).Scan(&request.ID)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" && strings.Contains(pqErr.Constraint, "reference") {
		return nil, ErrDuplicateCateringReference
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create catering request: %w", err)
	}

	return request, nil
}

func (r *CateringRepository) GetByID(id int) (*models.CateringRequest, error) {
	return r.getOne(cateringSelect+` WHERE cr.id = $1`, id)
}

// GetByReference finds a request by its customer-facing reference.
func (r *CateringRepository) GetByReference(reference string) (*models.CateringRequest, error) {
	return r.getOne(cateringSelect+` WHERE cr.reference = $1`, strings.ToUpper(strings.TrimSpace(reference)))
}

func (r *CateringRepository) GetByUserID(userID int) ([]models.CateringRequest, error) {
	return r.getMany(cateringSelect+` WHERE cr.user_id = $1 ORDER BY cr.created_at DESC`, userID)
}

func (r *CateringRepository) GetAll() ([]models.CateringRequest, error) {
	return r.getMany(cateringSelect + ` ORDER BY cr.created_at DESC`)
}

func (r *CateringRepository) UpdateStatus(id int, status string) error {
	result, err := r.db.Exec(
		`UPDATE catering_requests SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("failed to update catering request status: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ApplyUpdate saves the admin's notes; nil fields keep their current value.
// Status changes go through the contract repository so they are logged.
func (r *CateringRepository) ApplyUpdate(id int, update *models.CateringRequestUpdate) error {
	result, err := r.db.Exec(`
		UPDATE catering_requests SET
			admin_notes = COALESCE($1, admin_notes),
			updated_at = $2
		WHERE id = $3`,
		update.AdminNotes, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("failed to update catering request: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CateringRepository) getOne(query string, arg any) (*models.CateringRequest, error) {
	requests, err := r.getMany(query, arg)
	if err != nil || len(requests) == 0 {
		return nil, err
	}
	return &requests[0], nil
}

func (r *CateringRepository) getMany(query string, args ...any) ([]models.CateringRequest, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get catering requests: %w", err)
	}
	defer rows.Close()

	var requests []models.CateringRequest
	for rows.Next() {
		request, err := scanCateringRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, *request)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read catering requests: %w", err)
	}
	return requests, nil
}

func scanCateringRequest(row rowScanner) (*models.CateringRequest, error) {
	var request models.CateringRequest
	var (
		reference                  sql.NullString
		userID, packageID          sql.NullInt64
		eventDate                  time.Time
		eventTime                  sql.NullString
		duration                   sql.NullFloat64
		budget, estimated, quoted  sql.NullFloat64
		selectedItems, dietary     []byte
		updatedAt                  sql.NullTime
		acceptedQuoteID            sql.NullInt64
		contractedAt, paidInFullAt sql.NullTime
		eventType                  sql.NullString
	)

	err := row.Scan(
		&request.ID, &reference, &request.RequestType, &userID, &packageID,
		&request.PackageName,
		&request.ContactName, &request.ContactPhone, &request.ContactEmail, &request.Organization,
		&eventType, &request.EventTypeOther, &request.EventName, &eventDate,
		&eventTime, &duration, &request.GuestCount,
		&request.VenueName, &request.VenueAddress, &request.VenueArea, &request.VenueState, &request.VenueNotes,
		&request.ServiceStyle, &selectedItems, &request.MenuNotes,
		&dietary, &request.DietaryNotes, &budget, &request.SpecialRequests,
		&estimated, &quoted, &request.AdminNotes, &request.Status,
		&request.CreatedAt, &updatedAt, &acceptedQuoteID, &request.InvoiceNumber,
		&contractedAt, &paidInFullAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan catering request: %w", err)
	}

	request.Reference = reference.String
	request.EventType = eventType.String
	request.EventDate = eventDate.Format("2006-01-02")
	if eventTime.Valid {
		t := clockString(eventTime.String)
		request.EventTime = &t
	}
	request.UserID = nullIntPtr(userID)
	request.PackageID = nullIntPtr(packageID)
	request.DurationHours = nullFloatPtr(duration)
	request.Budget = nullFloatPtr(budget)
	request.EstimatedTotal = nullFloatPtr(estimated)
	request.QuotedAmount = nullFloatPtr(quoted)
	if updatedAt.Valid {
		request.UpdatedAt = &updatedAt.Time
	}
	request.AcceptedQuoteID = nullIntPtr(acceptedQuoteID)
	if contractedAt.Valid {
		request.ContractedAt = &contractedAt.Time
	}
	if paidInFullAt.Valid {
		request.PaidInFullAt = &paidInFullAt.Time
	}

	request.SelectedItems = []models.CateringSelectedItem{}
	if len(selectedItems) > 0 {
		if err := json.Unmarshal(selectedItems, &request.SelectedItems); err != nil {
			return nil, fmt.Errorf("failed to parse selected items: %w", err)
		}
	}
	request.DietaryRequirements = []string{}
	if len(dietary) > 0 {
		if err := json.Unmarshal(dietary, &request.DietaryRequirements); err != nil {
			return nil, fmt.Errorf("failed to parse dietary requirements: %w", err)
		}
	}

	return &request, nil
}

// clockString turns a TIME column into "15:04". lib/pq reads TIME as
// "0000-01-01T15:04:05Z"; other drivers return "15:04:05".
func clockString(value string) string {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.Format("15:04")
	}
	if len(value) >= 5 {
		return value[:5]
	}
	return value
}

func nullIntPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	value := int(v.Int64)
	return &value
}
