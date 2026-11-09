package repository

import (
	"context"
	"database/sql"
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// AdminPaymentRepository reads online payments for the admin.
type AdminPaymentRepository struct {
	db *sql.DB
}

func NewAdminPaymentRepository(db *sql.DB) *AdminPaymentRepository {
	return &AdminPaymentRepository{db: db}
}

const onlinePaymentWhere = `
	WHERE o.payment_method = 'transfer'
	  AND ($1 = '' OR o.payment_status = $1)
	  AND ($2 = '' OR o.order_number ILIKE $2 OR o.customer_name ILIKE $2 OR o.customer_phone ILIKE $2
	       OR EXISTS (SELECT 1 FROM payment_references pr WHERE pr.order_id = o.id AND pr.reference ILIKE $2))
`

// List returns online transfer orders matching filter, newest first, and how
// many match.
func (r *AdminPaymentRepository) List(ctx context.Context, filter models.OnlinePaymentFilter, limit, offset int) ([]models.OnlinePayment, int, error) {
	search := customerSearchPattern(filter.Search)

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders o`+onlinePaymentWhere,
		filter.PaymentStatus, search).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count online payments: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT o.id, o.order_number, o.customer_name, o.customer_phone, o.total_amount, o.status,
		       o.payment_status, o.payment_reference, o.payment_received_kobo, o.payment_review_reason,
		       (SELECT COUNT(*) FROM payment_references pr WHERE pr.order_id = o.id),
		       o.payment_expires_at, o.payment_paid_at, o.created_at
		FROM orders o`+onlinePaymentWhere+`
		ORDER BY o.created_at DESC, o.id DESC
		LIMIT $3 OFFSET $4
	`, filter.PaymentStatus, search, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list online payments: %w", err)
	}
	defer rows.Close()

	payments := []models.OnlinePayment{}
	for rows.Next() {
		var p models.OnlinePayment
		var reference, reviewReason sql.NullString
		var receivedKobo sql.NullInt64
		var expiresAt, paidAt sql.NullTime
		if err := rows.Scan(&p.OrderID, &p.OrderNumber, &p.CustomerName, &p.CustomerPhone, &p.TotalAmount, &p.Status,
			&p.PaymentStatus, &reference, &receivedKobo, &reviewReason, &p.ReferenceCount,
			&expiresAt, &paidAt, &p.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("read online payment: %w", err)
		}
		p.PaymentReference, p.ReviewReason = nullableString(reference), nullableString(reviewReason)
		if receivedKobo.Valid {
			amount := float64(receivedKobo.Int64) / 100
			p.ReceivedAmount = &amount
		}
		if expiresAt.Valid {
			p.ExpiresAt = &expiresAt.Time
		}
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		payments = append(payments, p)
	}
	return payments, total, rows.Err()
}

// CountByStatus returns how many online transfer orders have each payment
// status, for the page's tabs.
func (r *AdminPaymentRepository) CountByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT payment_status, COUNT(*) FROM orders WHERE payment_method = 'transfer' GROUP BY payment_status
	`)
	if err != nil {
		return nil, fmt.Errorf("count online payments by status: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{}
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			return nil, fmt.Errorf("read payment count: %w", err)
		}
		counts[status] = n
	}
	return counts, rows.Err()
}

// Unmatched returns payments Paystack reported that matched no order,
// newest first.
func (r *AdminPaymentRepository) Unmatched(ctx context.Context, limit int) ([]models.PaymentEventRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, event, reference, amount_kobo, currency, outcome, received_at
		FROM payment_events WHERE outcome = 'unmatched'
		ORDER BY received_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list unmatched payments: %w", err)
	}
	defer rows.Close()

	events := []models.PaymentEventRecord{}
	for rows.Next() {
		var e models.PaymentEventRecord
		var kobo sql.NullInt64
		var currency sql.NullString
		if err := rows.Scan(&e.ID, &e.Event, &e.Reference, &kobo, &currency, &e.Outcome, &e.ReceivedAt); err != nil {
			return nil, fmt.Errorf("read unmatched payment: %w", err)
		}
		if kobo.Valid {
			amount := float64(kobo.Int64) / 100
			e.Amount = &amount
		}
		e.Currency = nullableString(currency)
		events = append(events, e)
	}
	return events, rows.Err()
}
