package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// checkoutFixture creates a menu item with stock and a promotion, and removes
// everything the test created afterwards. Runs only with TEST_DATABASE_URL.
type checkoutFixture struct {
	db          *sql.DB
	repo        *CheckoutRepository
	idempotency *IdempotencyRepository
	menuItemID  int
	promotionID int
	tag         string
}

func newCheckoutFixture(t *testing.T, stock, promotionLimit int) *checkoutFixture {
	t.Helper()
	db := openIntegrationDB(t)
	f := &checkoutFixture{db: db, tag: fmt.Sprint(time.Now().UnixNano())}

	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(db.QueryRow(`INSERT INTO menu_items (category_id, name, price, is_available, created_at, updated_at)
		VALUES (NULL, $1, 1000, true, NOW(), NOW()) RETURNING id`, "Checkout test "+f.tag).Scan(&f.menuItemID))
	_, err := db.Exec(`INSERT INTO inventory_items (menu_item_id, name, current_stock, minimum_stock, reorder_point, unit, is_active, created_at, updated_at)
		VALUES ($1, 'Checkout test', $2, 0, 0, 'plate', true, NOW(), NOW())`, f.menuItemID, stock)
	must(err)
	must(db.QueryRow(`INSERT INTO promotions (code, title, promotion_type, discount_value, valid_from, valid_until, usage_limit, used_count, is_active, created_at, updated_at)
		VALUES ($1, 'Checkout test', 'fixed', 100, NOW() - INTERVAL '1 day', NOW() + INTERVAL '1 day', $2, 0, true, NOW(), NOW()) RETURNING id`,
		"CHK"+f.tag[len(f.tag)-8:], promotionLimit).Scan(&f.promotionID))

	t.Cleanup(func() {
		db.Exec(`DELETE FROM promotion_usage WHERE promotion_id = $1`, f.promotionID)
		db.Exec(`DELETE FROM order_items WHERE menu_item_id = $1`, f.menuItemID)
		db.Exec(`DELETE FROM audit_log WHERE entity_type = 'order' AND entity_id IN (SELECT id::text FROM orders WHERE order_number LIKE $1)`, "CHK-"+f.tag+"%")
		db.Exec(`DELETE FROM jobs WHERE payload->>'orderId' IN (SELECT id::text FROM orders WHERE order_number LIKE $1)`, "CHK-"+f.tag+"%")
		db.Exec(`DELETE FROM orders WHERE order_number LIKE $1`, "CHK-"+f.tag+"%")
		db.Exec(`DELETE FROM users WHERE full_name LIKE $1`, "Checkout guest "+f.tag+"%")
		db.Exec(`DELETE FROM inventory_alerts WHERE inventory_item_id IN (SELECT id FROM inventory_items WHERE menu_item_id = $1)`, f.menuItemID)
		db.Exec(`DELETE FROM inventory_history WHERE inventory_item_id IN (SELECT id FROM inventory_items WHERE menu_item_id = $1)`, f.menuItemID)
		db.Exec(`DELETE FROM inventory_items WHERE menu_item_id = $1`, f.menuItemID)
		db.Exec(`DELETE FROM menu_items WHERE id = $1`, f.menuItemID)
		db.Exec(`DELETE FROM promotions WHERE id = $1`, f.promotionID)
		db.Exec(`DELETE FROM idempotency_keys WHERE scope = $1`, "POST /test-checkout/"+f.tag)
	})

	orders := NewOrderRepository(db)
	f.idempotency = NewIdempotencyRepository(db)
	f.repo = NewCheckoutRepository(db, orders, NewUserRepository(db), NewInventoryRepository(db),
		NewPromotionRepository(db), NewAuditRepository(db), f.idempotency, NewJobRepository(db))
	return f
}

func (f *checkoutFixture) checkout(n int, quantity int, phone string, email *string, promotion bool) *models.Checkout {
	checkout := &models.Checkout{
		Order: &models.Order{
			OrderNumber:   fmt.Sprintf("CHK-%s-%d", f.tag, n),
			CustomerName:  "Checkout guest " + f.tag,
			CustomerPhone: phone,
			CustomerEmail: email,
			OrderType:     models.OrderTypePickup,
			Status:        models.OrderStatusPending,
			TotalAmount:   float64(quantity) * 1000,
			CreatedAt:     time.Now(),
			PaymentMethod: models.PaymentMethodCash,
			PaymentStatus: models.PaymentStatusPending,
		},
		Items: []models.OrderItemRequest{{MenuItemID: f.menuItemID, Name: "Checkout test", Quantity: quantity, UnitPrice: 1000}},
		Guest: &models.GuestCustomer{Name: "Checkout guest " + f.tag, Phone: phone, Email: email},
	}
	if promotion {
		checkout.Promotion = &models.PromotionUsage{PromotionID: f.promotionID, DiscountApplied: 100}
	}
	return checkout
}

func (f *checkoutFixture) count(t *testing.T, query string, args ...any) int {
	t.Helper()
	var n int
	if err := f.db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func (f *checkoutFixture) stock(t *testing.T) int {
	return f.count(t, `SELECT current_stock FROM inventory_items WHERE menu_item_id = $1`, f.menuItemID)
}

func TestPlaceOrderSavesEverythingTogether(t *testing.T) {
	f := newCheckoutFixture(t, 5, 10)
	phone := "0809" + f.tag[len(f.tag)-7:]
	email := "guest" + f.tag + "@example.test"

	scope := "POST /test-checkout/" + f.tag
	ctx := requestctx.WithIdempotency(context.Background(), requestctx.Idempotency{Scope: scope, Key: "checkout-key-1"})
	if outcome, _, err := f.idempotency.Begin(ctx, scope, "checkout-key-1", "hash", time.Minute, time.Hour); err != nil || outcome != models.IdempotencyNew {
		t.Fatalf("claim key: %v %v", outcome, err)
	}

	orderID, err := f.repo.PlaceOrder(ctx, f.checkout(1, 2, phone, &email, true))
	if err != nil {
		t.Fatal(err)
	}

	if got := f.stock(t); got != 3 {
		t.Fatalf("expected stock 3 after taking 2 of 5, got %d", got)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM orders o JOIN users u ON u.id = o.user_id WHERE o.id = $1 AND u.phone = $2 AND u.email = $3`, orderID, phone, email); n != 1 {
		t.Fatal("expected the order linked to a new customer record with the guest's phone and email")
	}
	if n := f.count(t, `SELECT COUNT(*) FROM order_items WHERE order_id = $1`, orderID); n != 1 {
		t.Fatalf("expected 1 order item, got %d", n)
	}
	if n := f.count(t, `SELECT used_count FROM promotions WHERE id = $1`, f.promotionID); n != 1 {
		t.Fatalf("expected the promotion used once, got %d", n)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM promotion_usage WHERE order_id = $1 AND customer_id IS NOT NULL`, orderID); n != 1 {
		t.Fatal("expected the promotion use recorded against the order and customer")
	}
	if n := f.count(t, `SELECT COUNT(*) FROM audit_log WHERE action = 'order.created' AND entity_id = $1 AND idempotency_key = 'checkout-key-1'`, fmt.Sprint(orderID)); n != 1 {
		t.Fatal("expected an audit entry carrying the idempotency key")
	}
	if outcome, record, _ := f.idempotency.Begin(ctx, scope, "checkout-key-1", "hash", time.Minute, time.Hour); outcome != models.IdempotencyReplay || record.EntityID != fmt.Sprint(orderID) {
		t.Fatalf("expected the key completed with the order, got %v %+v", outcome, record)
	}
}

func TestPlaceOrderSavesNothingWhenStockRunsShort(t *testing.T) {
	f := newCheckoutFixture(t, 1, 10)
	phone := "0808" + f.tag[len(f.tag)-7:]

	_, err := f.repo.PlaceOrder(context.Background(), f.checkout(1, 2, phone, nil, true))
	var shortage *InsufficientStockError
	if !errors.As(err, &shortage) || shortage.Available != 1 {
		t.Fatalf("expected a shortage with 1 available, got %v", err)
	}

	if got := f.stock(t); got != 1 {
		t.Fatalf("stock must be unchanged, got %d", got)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM users WHERE phone = $1`, phone); n != 0 {
		t.Fatal("no customer record may remain from a rejected checkout")
	}
	if n := f.count(t, `SELECT COUNT(*) FROM orders WHERE order_number LIKE $1`, "CHK-"+f.tag+"%"); n != 0 {
		t.Fatal("no order may remain from a rejected checkout")
	}
	if n := f.count(t, `SELECT used_count FROM promotions WHERE id = $1`, f.promotionID); n != 0 {
		t.Fatal("the promotion must not be used by a rejected checkout")
	}
}

func TestPlaceOrderSavesNothingWhenPromotionRunsOut(t *testing.T) {
	f := newCheckoutFixture(t, 5, 0)

	_, err := f.repo.PlaceOrder(context.Background(), f.checkout(1, 1, "0807"+f.tag[len(f.tag)-7:], nil, true))
	if !errors.Is(err, ErrPromotionUnavailable) {
		t.Fatalf("expected ErrPromotionUnavailable, got %v", err)
	}
	if got := f.stock(t); got != 5 {
		t.Fatalf("stock taken before the promotion failed must be returned, got %d", got)
	}
}

func TestPlaceOrderNeverOversellsUnderConcurrency(t *testing.T) {
	f := newCheckoutFixture(t, 1, 10)

	const buyers = 6
	var wg sync.WaitGroup
	results := make(chan error, buyers)
	for i := range buyers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_, err := f.repo.PlaceOrder(context.Background(), f.checkout(n, 1, fmt.Sprintf("0806%s%d", f.tag[len(f.tag)-6:], n), nil, false))
			results <- err
		}(i)
	}
	wg.Wait()
	close(results)

	placed := 0
	for err := range results {
		var shortage *InsufficientStockError
		switch {
		case err == nil:
			placed++
		case !errors.As(err, &shortage):
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if placed != 1 || f.stock(t) != 0 {
		t.Fatalf("expected exactly 1 of %d buyers to get the last item and stock 0, got %d placed and stock %d", buyers, placed, f.stock(t))
	}
}

func TestPlaceOrderRecordsGuestWhoseEmailBelongsToAnotherCustomer(t *testing.T) {
	f := newCheckoutFixture(t, 5, 10)
	email := "shared" + f.tag + "@example.test"
	first := "0805" + f.tag[len(f.tag)-7:]
	second := "0804" + f.tag[len(f.tag)-7:]

	if _, err := f.repo.PlaceOrder(context.Background(), f.checkout(1, 1, first, &email, false)); err != nil {
		t.Fatal(err)
	}
	orderID, err := f.repo.PlaceOrder(context.Background(), f.checkout(2, 1, second, &email, false))
	if err != nil {
		t.Fatal(err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM orders o JOIN users u ON u.id = o.user_id WHERE o.id = $1 AND u.phone = $2 AND u.email IS NULL`, orderID, second); n != 1 {
		t.Fatal("expected the second guest recorded by phone without the taken email")
	}

	// A returning guest is matched by phone, not duplicated.
	if _, err := f.repo.PlaceOrder(context.Background(), f.checkout(3, 1, first, &email, false)); err != nil {
		t.Fatal(err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM users WHERE phone = $1`, first); n != 1 {
		t.Fatalf("expected one record for the returning guest, got %d", n)
	}
}

func TestPlaceOrderQueuesFollowUpJobs(t *testing.T) {
	f := newCheckoutFixture(t, 5, 10)

	cash := f.checkout(1, 1, "0803"+f.tag[len(f.tag)-7:], nil, false)
	invoice := "INV-" + f.tag
	cash.Order.InvoiceNumber = &invoice
	cashID, err := f.repo.PlaceOrder(context.Background(), cash)
	if err != nil {
		t.Fatal(err)
	}
	transfer := f.checkout(2, 1, "0802"+f.tag[len(f.tag)-7:], nil, false)
	transfer.Order.PaymentMethod = models.PaymentMethodTransfer
	transferID, err := f.repo.PlaceOrder(context.Background(), transfer)
	if err != nil {
		t.Fatal(err)
	}

	jobsFor := func(orderID int) map[string]time.Time {
		rows, err := f.db.Query(`SELECT kind, run_at FROM jobs WHERE payload->>'orderId' = $1`, fmt.Sprint(orderID))
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		kinds := map[string]time.Time{}
		for rows.Next() {
			var kind string
			var runAt time.Time
			_ = rows.Scan(&kind, &runAt)
			kinds[kind] = runAt
		}
		return kinds
	}

	if got := jobsFor(cashID); len(got) != 1 || got["order.invoice"].IsZero() {
		t.Fatalf("a pay-on-delivery order queues only its invoice email, got %v", got)
	}
	got := jobsFor(transferID)
	if len(got) != 2 || got["order.confirmation"].IsZero() {
		t.Fatalf("a transfer order queues its confirmation and a payment check, got %v", got)
	}
	if check := got["order.payment_check"]; time.Until(check) < PaymentCheckDelay-time.Minute {
		t.Fatalf("the payment check must wait for the window to close, runs at %v", check)
	}
}

func TestCancelUnpaidOrderReturnsStockOnce(t *testing.T) {
	f := newCheckoutFixture(t, 5, 10)
	transfer := f.checkout(1, 2, "0801"+f.tag[len(f.tag)-7:], nil, false)
	transfer.Order.PaymentMethod = models.PaymentMethodTransfer
	orderID, err := f.repo.PlaceOrder(context.Background(), transfer)
	if err != nil {
		t.Fatal(err)
	}
	if f.stock(t) != 3 {
		t.Fatalf("expected stock 3 after the order, got %d", f.stock(t))
	}

	cancelled, err := f.repo.CancelUnpaidOrder(context.Background(), orderID)
	if err != nil || !cancelled {
		t.Fatalf("expected the unpaid order cancelled, got %v %v", cancelled, err)
	}
	if f.stock(t) != 5 {
		t.Fatalf("expected the stock returned, got %d", f.stock(t))
	}
	if n := f.count(t, `SELECT COUNT(*) FROM orders WHERE id = $1 AND status = 'cancelled' AND payment_status = 'expired'`, orderID); n != 1 {
		t.Fatal("expected the order cancelled with its payment expired")
	}
	if n := f.count(t, `SELECT COUNT(*) FROM audit_log WHERE action = 'order.payment_expired' AND entity_id = $1`, fmt.Sprint(orderID)); n != 1 {
		t.Fatal("expected an audit entry for the expiry")
	}

	// Running the check again changes nothing.
	if again, err := f.repo.CancelUnpaidOrder(context.Background(), orderID); err != nil || again || f.stock(t) != 5 {
		t.Fatalf("a second cancellation must do nothing, got %v %v stock %d", again, err, f.stock(t))
	}

	// A paid order is never cancelled.
	paid := f.checkout(2, 1, "0801"+f.tag[len(f.tag)-7:], nil, false)
	paid.Order.PaymentMethod = models.PaymentMethodTransfer
	paidID, _ := f.repo.PlaceOrder(context.Background(), paid)
	_, _ = f.db.Exec(`UPDATE orders SET payment_status = 'paid' WHERE id = $1`, paidID)
	if c, err := f.repo.CancelUnpaidOrder(context.Background(), paidID); err != nil || c {
		t.Fatalf("a paid order must not be cancelled, got %v %v", c, err)
	}
}
