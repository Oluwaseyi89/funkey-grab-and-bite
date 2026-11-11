package repository

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Runs only with TEST_DATABASE_URL. Orders are dated in 2001 so other data
// in the database doesn't affect the figures.
func TestReportsCountPaidRevenueWithEveryDay(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()
	repo := NewReportsRepository(db)

	// Use the session's own zone, so stored times read back unchanged.
	var tz string
	if err := db.QueryRow(`SELECT current_setting('TimeZone')`).Scan(&tz); err != nil {
		t.Fatal(err)
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		t.Skipf("session timezone %q isn't a loadable zone: %v", tz, err)
	}

	tag := fmt.Sprint(time.Now().UnixNano())
	insert := func(day, status, payment string, amount float64) {
		t.Helper()
		_, err := db.Exec(`
			INSERT INTO orders (order_number, customer_name, customer_phone, order_type, status, total_amount, payment_method, payment_status, created_at, updated_at)
			VALUES ($1, 'Report test', '08000000000', 'pickup', $2, $3, 'cash', $4, $5::timestamp, now())`,
			fmt.Sprintf("RPT-%s-%d", tag, time.Now().UnixNano()), status, amount, payment, day+" 12:00:00")
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM orders WHERE order_number LIKE $1`, "RPT-"+tag+"-%") })

	insert("2001-03-05", "completed", "paid", 1000) // Monday
	insert("2001-03-05", "cancelled", "paid", 9000) // cancelled: no revenue, no order
	insert("2001-03-05", "pending", "pending", 500) // placed, not paid
	insert("2001-03-07", "completed", "paid", 2000) // Wednesday
	insert("2001-03-11", "completed", "paid", 4000) // Sunday, the last day

	day := func(s string) time.Time { d, _ := time.ParseInLocation("2006-01-02", s, loc); return d }

	totals, err := repo.Totals(ctx, tz, day("2001-03-05"), day("2001-03-12"))
	if err != nil {
		t.Fatal(err)
	}
	if totals.Revenue != 7000 || totals.PaidOrders != 3 || totals.Orders != 4 || totals.Cancelled != 1 {
		t.Fatalf("expected ₦7000 from 3 paid of 4 orders and 1 cancelled, got %+v", totals)
	}

	series, err := repo.Series(ctx, tz, day("2001-03-05"), day("2001-03-11"), "day")
	if err != nil {
		t.Fatal(err)
	}
	if len(series) != 7 {
		t.Fatalf("expected all 7 days including empty ones, got %d: %+v", len(series), series)
	}
	if series[1].Date != "2001-03-06" || series[1].Revenue != 0 || series[1].Orders != 0 {
		t.Fatalf("a day with no sales must be zero, got %+v", series[1])
	}
	if last := series[6]; last.Date != "2001-03-11" || last.Revenue != 4000 {
		t.Fatalf("the last day must be included, got %+v", last)
	}

	weeks, err := repo.Series(ctx, tz, day("2001-03-05"), day("2001-03-11"), "week")
	if err != nil {
		t.Fatal(err)
	}
	if len(weeks) != 1 || weeks[0].Date != "2001-03-05" || weeks[0].Revenue != 7000 || weeks[0].Orders != 4 {
		t.Fatalf("expected one week starting Monday with everything added up, got %+v", weeks)
	}

	hours, err := repo.Hourly(ctx, tz, day("2001-03-05"), day("2001-03-12"))
	if err != nil {
		t.Fatal(err)
	}
	if len(hours) != 24 || hours[12].Orders != 4 {
		t.Fatalf("expected 4 orders at noon across 24 hours, got %d hours, noon %+v", len(hours), hours[12])
	}
}
