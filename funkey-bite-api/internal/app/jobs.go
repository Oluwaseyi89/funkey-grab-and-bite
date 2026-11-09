package app

import (
	"context"
	"encoding/json"
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/jobs"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// registerJobs tells the worker how to run each kind of job.
func registerJobs(worker *jobs.Worker, orders services.OrderService, notifications services.NotificationService, invoices services.InvoiceService, payments services.PaymentService, email utils.EmailService) {
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

	// The payment-confirmed messages, queued when a transfer is recorded.
	worker.Register(jobs.KindOrderPaidEmail, func(ctx context.Context, payload json.RawMessage) error {
		orderID, err := jobs.Order(payload)
		if err != nil {
			return err
		}
		return notifications.SendOrderStatusEmail(orderID, string(models.OrderStatusConfirmed))
	})
	worker.Register(jobs.KindOrderPaidSMS, func(ctx context.Context, payload json.RawMessage) error {
		orderID, err := jobs.Order(payload)
		if err != nil {
			return err
		}
		return notifications.SendOrderStatusSMS(orderID, string(models.OrderStatusConfirmed))
	})

	worker.Register(jobs.KindEmail, func(ctx context.Context, payload json.RawMessage) error {
		var message models.EmailJobPayload
		if err := json.Unmarshal(payload, &message); err != nil || message.To == "" {
			return jobs.Permanent(fmt.Errorf("invalid email job payload"))
		}
		return email.SendInvoice(message.To, message.Subject, message.HTML, message.Text)
	})

	worker.Register(jobs.KindOrderPaymentCheck, func(ctx context.Context, payload json.RawMessage) error {
		orderID, err := jobs.Order(payload)
		if err != nil {
			return err
		}
		return payments.CheckTransferPayment(ctx, orderID)
	})
}

// jobEmailQueue queues rendered emails as sensitive jobs and wakes the
// worker, so they are sent at once and retried if sending fails.
type jobEmailQueue struct {
	jobs   *repository.JobRepository
	worker *jobs.Worker
}

func (q jobEmailQueue) QueueEmail(ctx context.Context, email models.EmailJobPayload) error {
	if err := q.jobs.Enqueue(ctx, nil, models.NewJob{Kind: jobs.KindEmail, Payload: email, Sensitive: true}); err != nil {
		return err
	}
	q.worker.Notify()
	return nil
}
