package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/database"
	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/jobs"
)

// CheckoutRepository saves an order and everything it changes in a single
// transaction.
type CheckoutRepository struct {
	db          *sql.DB
	orders      *OrderRepository
	users       *UserRepository
	inventory   *InventoryRepository
	promotions  *PromotionRepository
	audit       *AuditRepository
	idempotency *IdempotencyRepository
	jobs        *JobRepository
}

// PaymentCheckDelay is how long after checkout an online transfer order is
// checked for payment: the 15-minute payment window plus time for a late
// transfer to be confirmed. The check waits longer if the window was renewed.
const PaymentCheckDelay = 20 * time.Minute

func NewCheckoutRepository(db *sql.DB, orders *OrderRepository, users *UserRepository, inventory *InventoryRepository, promotions *PromotionRepository, audit *AuditRepository, idempotency *IdempotencyRepository, jobs *JobRepository) *CheckoutRepository {
	return &CheckoutRepository{
		db: db, orders: orders, users: users, inventory: inventory,
		promotions: promotions, audit: audit, idempotency: idempotency, jobs: jobs,
	}
}

// PlaceOrder saves checkout and returns the new order's ID. In one
// transaction it finds or creates the guest's customer record, inserts the
// order and its items, takes the stock (locking each row, and failing with an
// *InsufficientStockError rather than going below zero), records the
// promotion use (failing with ErrPromotionUnavailable if it ran out), writes
// the audit entry, queues the follow-up jobs (confirmation or invoice email,
// and for online transfers a payment check) and completes the request's
// idempotency key. Either all of it is saved or none of it is.
func (r *CheckoutRepository) PlaceOrder(ctx context.Context, checkout *models.Checkout) (int, error) {
	var orderID int
	err := database.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		// Work on a copy: the transaction may be retried from the start.
		order := *checkout.Order
		if order.UserID == nil && checkout.Guest != nil {
			guest := checkout.Guest
			customerID, err := r.users.FindOrCreateGuestTx(ctx, tx, guest.Name, guest.Phone, guest.Email)
			if err != nil {
				return err
			}
			order.UserID = &customerID
		}

		created, err := r.orders.CreateOrderWithTransaction(tx, &order)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		for _, item := range checkout.Items {
			_, err := r.orders.CreateOrderItemWithTransaction(tx, &models.OrderItem{
				OrderID:             created.ID,
				MenuItemID:          item.MenuItemID,
				Name:                item.Name,
				Quantity:            item.Quantity,
				UnitPrice:           item.UnitPrice,
				SpecialInstructions: item.SpecialInstructions,
			})
			if err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		// Lock stock rows in menu item order so two checkouts sharing items
		// can't deadlock waiting for each other.
		items := slices.Clone(checkout.Items)
		slices.SortFunc(items, func(a, b models.OrderItemRequest) int { return a.MenuItemID - b.MenuItemID })
		reason := "Order #" + order.OrderNumber
		for _, item := range items {
			inventoryID, err := r.inventory.AdjustStockTx(tx, item.MenuItemID, item.Quantity, "subtract", reason)
			if err != nil {
				return err
			}
			if inventoryID == 0 {
				return &InsufficientStockError{MenuItemID: item.MenuItemID, Requested: item.Quantity, Untracked: true}
			}
		}

		if checkout.Promotion != nil {
			usage := *checkout.Promotion
			usage.OrderID = created.ID
			usage.CustomerID = order.UserID
			if err := r.promotions.ConsumeUsageTx(tx, &usage); err != nil {
				return err
			}
		}

		entityID := strconv.Itoa(created.ID)
		if err := r.audit.Record(ctx, tx, models.AuditEntry{
			Action:     "order.created",
			EntityType: "order",
			EntityID:   entityID,
			After: map[string]any{
				"orderNumber":   order.OrderNumber,
				"orderType":     order.OrderType,
				"status":        order.Status,
				"paymentMethod": order.PaymentMethod,
				"totalAmount":   order.TotalAmount,
				"itemCount":     len(checkout.Items),
				"customerId":    order.UserID,
			},
		}); err != nil {
			return err
		}
		if err := r.queueFollowUps(ctx, tx, created.ID, &order); err != nil {
			return err
		}
		if err := r.idempotency.MarkCompleted(ctx, tx, "order", entityID); err != nil {
			return err
		}

		orderID = created.ID
		return nil
	})
	return orderID, err
}

// queueFollowUps queues what happens after an order is placed. The jobs are
// part of the order's transaction, so they exist only if the order does.
func (r *CheckoutRepository) queueFollowUps(ctx context.Context, tx *sql.Tx, orderID int, order *models.Order) error {
	payload := models.OrderJobPayload{OrderID: orderID}
	id := strconv.Itoa(orderID)

	// Pay-on-delivery orders get the invoice, which covers the confirmation.
	email := models.NewJob{Kind: jobs.KindOrderConfirmation, Payload: payload, DedupeKey: jobs.KindOrderConfirmation + ":" + id}
	if order.InvoiceNumber != nil {
		email = models.NewJob{Kind: jobs.KindOrderInvoice, Payload: payload, DedupeKey: jobs.KindOrderInvoice + ":" + id}
	}
	if err := r.jobs.Enqueue(ctx, tx, email); err != nil {
		return err
	}

	if order.PaymentMethod == models.PaymentMethodTransfer {
		return r.jobs.Enqueue(ctx, tx, models.NewJob{
			Kind:      jobs.KindOrderPaymentCheck,
			Payload:   payload,
			RunAt:     time.Now().Add(PaymentCheckDelay),
			DedupeKey: jobs.KindOrderPaymentCheck + ":" + id,
		})
	}
	return nil
}

// CancelUnpaidOrder cancels an online transfer order that was never paid and
// returns its stock, in one transaction with its audit entry. It returns false
// without changing anything when the order no longer qualifies, e.g. it was
// paid or cancelled in the meantime.
func (r *CheckoutRepository) CancelUnpaidOrder(ctx context.Context, orderID int) (bool, error) {
	cancelled := false
	err := database.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		cancelled = false
		var orderNumber, status, paymentStatus, paymentMethod string
		err := tx.QueryRowContext(ctx, `
			SELECT order_number, status, payment_status, payment_method
			FROM orders WHERE id = $1 FOR UPDATE
		`, orderID).Scan(&orderNumber, &status, &paymentStatus, &paymentMethod)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("lock order %d: %w", orderID, err)
		}
		if paymentMethod != string(models.PaymentMethodTransfer) ||
			paymentStatus != string(models.PaymentStatusPending) ||
			status == string(models.OrderStatusCancelled) {
			return nil
		}

		rows, err := tx.QueryContext(ctx, `
			SELECT menu_item_id, quantity FROM order_items WHERE order_id = $1 ORDER BY menu_item_id
		`, orderID)
		if err != nil {
			return fmt.Errorf("read order %d items: %w", orderID, err)
		}
		type line struct{ menuItemID, quantity int }
		var lines []line
		for rows.Next() {
			var l line
			if err := rows.Scan(&l.menuItemID, &l.quantity); err != nil {
				rows.Close()
				return fmt.Errorf("read order %d items: %w", orderID, err)
			}
			lines = append(lines, l)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE orders SET status = $2, payment_status = $3, updated_at = NOW() WHERE id = $1
		`, orderID, models.OrderStatusCancelled, models.PaymentStatusExpired); err != nil {
			return fmt.Errorf("cancel order %d: %w", orderID, err)
		}

		reason := "Unpaid order #" + orderNumber + " expired"
		for _, l := range lines {
			// An item no longer tracked in stock has nothing to return.
			if _, err := r.inventory.AdjustStockTx(tx, l.menuItemID, l.quantity, "add", reason); err != nil {
				return err
			}
		}

		if err := r.audit.Record(ctx, tx, models.AuditEntry{
			Action:     "order.payment_expired",
			EntityType: "order",
			EntityID:   strconv.Itoa(orderID),
			Before:     map[string]any{"status": status, "paymentStatus": paymentStatus},
			After:      map[string]any{"status": models.OrderStatusCancelled, "paymentStatus": models.PaymentStatusExpired},
			Metadata:   map[string]any{"stockReturned": len(lines) > 0},
		}); err != nil {
			return err
		}

		cancelled = true
		return nil
	})
	return cancelled, err
}
