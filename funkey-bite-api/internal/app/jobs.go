package app

import (
	"context"
	"encoding/json"
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/jobs"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

// registerJobs tells the worker how to run each kind of job.
func registerJobs(worker *jobs.Worker, orders services.OrderService, notifications services.NotificationService, invoices services.InvoiceService, payments services.PaymentService) {
	loadOrder := func(payload json.RawMessage) (*models.Order, error) {
		orderID, err := jobs.Order(payload)
		if err != nil {
			return nil, err
		}
		order, err := orders.GetOrderByID(orderID)
		if err != nil {
			return nil, fmt.Errorf("load order %d: %w", orderID, err)
		}
		if order == nil {
			return nil, jobs.Permanent(fmt.Errorf("order %d not found", orderID))
		}
		return order, nil
	}

	worker.Register(jobs.KindOrderConfirmation, func(ctx context.Context, payload json.RawMessage) error {
		order, err := loadOrder(payload)
		if err != nil {
			return err
		}
		return notifications.SendOrderConfirmation(order)
	})

	worker.Register(jobs.KindOrderInvoice, func(ctx context.Context, payload json.RawMessage) error {
		order, err := loadOrder(payload)
		if err != nil {
			return err
		}
		return invoices.SendInvoice(order)
	})

	worker.Register(jobs.KindOrderPaymentCheck, func(ctx context.Context, payload json.RawMessage) error {
		orderID, err := jobs.Order(payload)
		if err != nil {
			return err
		}
		return payments.CheckTransferPayment(ctx, orderID)
	})
}
