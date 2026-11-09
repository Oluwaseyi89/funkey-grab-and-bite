package repository

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strconv"

	"funkey-grab-and-bite/funkey-bite-api/internal/database"
	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
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
}

func NewCheckoutRepository(db *sql.DB, orders *OrderRepository, users *UserRepository, inventory *InventoryRepository, promotions *PromotionRepository, audit *AuditRepository, idempotency *IdempotencyRepository) *CheckoutRepository {
	return &CheckoutRepository{
		db: db, orders: orders, users: users, inventory: inventory,
		promotions: promotions, audit: audit, idempotency: idempotency,
	}
}

// PlaceOrder saves checkout and returns the new order's ID. In one
// transaction it finds or creates the guest's customer record, inserts the
// order and its items, takes the stock (locking each row, and failing with an
// *InsufficientStockError rather than going below zero), records the
// promotion use (failing with ErrPromotionUnavailable if it ran out), writes
// the audit entry and completes the request's idempotency key. Either all of
// it is saved or none of it is.
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
		if err := r.idempotency.MarkCompleted(ctx, tx, "order", entityID); err != nil {
			return err
		}

		orderID = created.ID
		return nil
	})
	return orderID, err
}
