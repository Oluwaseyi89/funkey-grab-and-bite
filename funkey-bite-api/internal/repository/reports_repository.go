package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// ReportsRepository works out the admin's sales figures. See models/reports.go
// for what revenue and orders mean.
type ReportsRepository struct {
	db *sql.DB
}

func NewReportsRepository(db *sql.DB) *ReportsRepository {
	return &ReportsRepository{db: db}
}

// Timestamps are stored without a zone, in the database session's zone.
// localTime turns a column into wall-clock time in the business timezone,
// which is always the first query parameter.
func localTime(column string) string {
	return fmt.Sprintf("((%s) AT TIME ZONE current_setting('TimeZone')) AT TIME ZONE $1", column)
}

const wallClock = "2006-01-02 15:04:05"

var (
	placedAt = localTime("o.created_at")
	// A payment counts on the day it was made; older rows without the time
	// fall back to when the order was placed.
	paidAt = localTime("COALESCE(o.payment_paid_at, o.created_at)")
)

const paidCondition = "o.payment_status = 'paid' AND o.status <> 'cancelled'"

// Totals sums revenue and orders between two wall-clock times, from
// included, to excluded.
func (r *ReportsRepository) Totals(ctx context.Context, tz string, from, to time.Time) (models.PeriodTotals, error) {
	var t models.PeriodTotals
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) FILTER (WHERE o.status <> 'cancelled' AND %[1]s >= $2::timestamp AND %[1]s < $3::timestamp),
			COUNT(*) FILTER (WHERE o.status = 'cancelled' AND %[1]s >= $2::timestamp AND %[1]s < $3::timestamp),
			COUNT(*) FILTER (WHERE %[3]s AND %[2]s >= $2::timestamp AND %[2]s < $3::timestamp),
			COALESCE(SUM(o.total_amount) FILTER (WHERE %[3]s AND %[2]s >= $2::timestamp AND %[2]s < $3::timestamp), 0)
		FROM orders o`, placedAt, paidAt, paidCondition)
	err := r.db.QueryRowContext(ctx, query, tz, from.Format(wallClock), to.Format(wallClock)).
		Scan(&t.Orders, &t.Cancelled, &t.PaidOrders, &t.Revenue)
	if err != nil {
		return t, fmt.Errorf("sales totals: %w", err)
	}
	if t.PaidOrders > 0 {
		t.AverageOrder = t.Revenue / float64(t.PaidOrders)
	}
	return t, nil
}

// Series returns one point per day, week or month from the period holding
// from to the period holding to (both dates included), with zeros where
// nothing was sold. group must be day, week or month.
func (r *ReportsRepository) Series(ctx context.Context, tz string, from, to time.Time, group string) ([]models.SalesPoint, error) {
	query := fmt.Sprintf(`
		WITH buckets AS (
			SELECT generate_series(date_trunc($4, $2::timestamp), date_trunc($4, $3::timestamp), ('1 ' || $4)::interval)::date AS bucket
		),
		placed AS (
			SELECT date_trunc($4, %[1]s)::date AS bucket, COUNT(*) AS orders
			FROM orders o
			WHERE o.status <> 'cancelled' AND %[1]s >= $2::timestamp AND %[1]s < ($3::timestamp + interval '1 day')
			GROUP BY 1
		),
		paid AS (
			SELECT date_trunc($4, %[2]s)::date AS bucket, COUNT(*) AS paid_orders, SUM(o.total_amount) AS revenue
			FROM orders o
			WHERE %[3]s AND %[2]s >= $2::timestamp AND %[2]s < ($3::timestamp + interval '1 day')
			GROUP BY 1
		)
		SELECT to_char(b.bucket, 'YYYY-MM-DD'), COALESCE(p.orders, 0), COALESCE(q.paid_orders, 0), COALESCE(q.revenue, 0)
		FROM buckets b
		LEFT JOIN placed p ON p.bucket = b.bucket
		LEFT JOIN paid q ON q.bucket = b.bucket
		ORDER BY b.bucket`, placedAt, paidAt, paidCondition)

	rows, err := r.db.QueryContext(ctx, query, tz, from.Format("2006-01-02"), to.Format("2006-01-02"), group)
	if err != nil {
		return nil, fmt.Errorf("sales series: %w", err)
	}
	defer rows.Close()

	points := []models.SalesPoint{}
	for rows.Next() {
		var p models.SalesPoint
		if err := rows.Scan(&p.Date, &p.Orders, &p.PaidOrders, &p.Revenue); err != nil {
			return nil, fmt.Errorf("scan sales point: %w", err)
		}
		points = append(points, p)
	}
	return points, rows.Err()
}

// Hourly counts orders placed in each hour of the day, all 24 hours included.
func (r *ReportsRepository) Hourly(ctx context.Context, tz string, from, to time.Time) ([]models.HourPoint, error) {
	query := fmt.Sprintf(`
		SELECT EXTRACT(HOUR FROM %[1]s)::int, COUNT(*)
		FROM orders o
		WHERE o.status <> 'cancelled' AND %[1]s >= $2::timestamp AND %[1]s < $3::timestamp
		GROUP BY 1`, placedAt)
	rows, err := r.db.QueryContext(ctx, query, tz, from.Format(wallClock), to.Format(wallClock))
	if err != nil {
		return nil, fmt.Errorf("orders by hour: %w", err)
	}
	defer rows.Close()

	hours := make([]models.HourPoint, 24)
	for i := range hours {
		hours[i].Hour = i
	}
	for rows.Next() {
		var hour, count int
		if err := rows.Scan(&hour, &count); err != nil {
			return nil, fmt.Errorf("scan hour: %w", err)
		}
		if hour >= 0 && hour < 24 {
			hours[hour].Orders = count
		}
	}
	return hours, rows.Err()
}

// TopItems lists the best sellers by quantity among orders placed in the
// period and not cancelled.
func (r *ReportsRepository) TopItems(ctx context.Context, tz string, from, to time.Time, limit int) ([]models.TopItem, error) {
	query := fmt.Sprintf(`
		SELECT oi.menu_item_id, MAX(oi.name), SUM(oi.quantity)::int, COALESCE(SUM(oi.quantity * oi.unit_price), 0)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE o.status <> 'cancelled' AND %[1]s >= $2::timestamp AND %[1]s < $3::timestamp
		GROUP BY oi.menu_item_id
		ORDER BY 3 DESC, 4 DESC
		LIMIT $4`, placedAt)
	rows, err := r.db.QueryContext(ctx, query, tz, from.Format(wallClock), to.Format(wallClock), limit)
	if err != nil {
		return nil, fmt.Errorf("top items: %w", err)
	}
	defer rows.Close()

	items := []models.TopItem{}
	for rows.Next() {
		var item models.TopItem
		if err := rows.Scan(&item.MenuItemID, &item.Name, &item.Quantity, &item.Revenue); err != nil {
			return nil, fmt.Errorf("scan top item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// PaymentMethods splits revenue received by how it was paid.
func (r *ReportsRepository) PaymentMethods(ctx context.Context, tz string, from, to time.Time) ([]models.Share, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(o.payment_method, 'unknown'), COUNT(*), COALESCE(SUM(o.total_amount), 0)
		FROM orders o
		WHERE %[2]s AND %[1]s >= $2::timestamp AND %[1]s < $3::timestamp
		GROUP BY 1
		ORDER BY 3 DESC`, paidAt, paidCondition)
	return r.shares(ctx, query, tz, from, to)
}

// OrderTypes splits orders placed by pickup and delivery.
func (r *ReportsRepository) OrderTypes(ctx context.Context, tz string, from, to time.Time) ([]models.Share, error) {
	query := fmt.Sprintf(`
		SELECT o.order_type, COUNT(*), COALESCE(SUM(o.total_amount), 0)
		FROM orders o
		WHERE o.status <> 'cancelled' AND %[1]s >= $2::timestamp AND %[1]s < $3::timestamp
		GROUP BY 1
		ORDER BY 2 DESC`, placedAt)
	return r.shares(ctx, query, tz, from, to)
}

func (r *ReportsRepository) shares(ctx context.Context, query, tz string, from, to time.Time) ([]models.Share, error) {
	rows, err := r.db.QueryContext(ctx, query, tz, from.Format(wallClock), to.Format(wallClock))
	if err != nil {
		return nil, fmt.Errorf("sales split: %w", err)
	}
	defer rows.Close()
	shares := []models.Share{}
	for rows.Next() {
		var s models.Share
		if err := rows.Scan(&s.Key, &s.Orders, &s.Revenue); err != nil {
			return nil, fmt.Errorf("scan sales split: %w", err)
		}
		shares = append(shares, s)
	}
	return shares, rows.Err()
}

// Workload counts what needs attention right now: orders in the kitchen,
// payments to check or still awaited, low stock and catering to quote.
func (r *ReportsRepository) Workload(ctx context.Context) (queue models.OrderQueue, toReview, awaiting, lowStock, cateringToQuote int, err error) {
	err = r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status = 'confirmed'),
			COUNT(*) FILTER (WHERE status = 'preparing'),
			COUNT(*) FILTER (WHERE status = 'ready'),
			COUNT(*) FILTER (WHERE payment_status = 'needs_review'),
			COUNT(*) FILTER (WHERE payment_status = 'pending' AND payment_method IN ('transfer', 'card') AND status <> 'cancelled')
		FROM orders`).
		Scan(&queue.Pending, &queue.Confirmed, &queue.Preparing, &queue.Ready, &toReview, &awaiting)
	if err != nil {
		err = fmt.Errorf("order workload: %w", err)
		return
	}
	if err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM inventory_items WHERE is_active = true AND current_stock <= reorder_point`).Scan(&lowStock); err != nil {
		err = fmt.Errorf("low stock count: %w", err)
		return
	}
	if err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM catering_requests WHERE status = 'pending'`).Scan(&cateringToQuote); err != nil {
		err = fmt.Errorf("catering to quote: %w", err)
	}
	return
}
