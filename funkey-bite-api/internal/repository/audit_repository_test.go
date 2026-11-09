package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

func TestAuditRecordFillsActorAndKeyFromContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := requestctx.WithActor(context.Background(), requestctx.Actor{
		Type: requestctx.ActorAdmin, ID: "3", IP: "198.51.100.4", RequestID: "req-1",
	})
	ctx = requestctx.WithIdempotency(ctx, requestctx.Idempotency{Scope: "POST /orders", Key: "checkout-abc-123"})

	mock.ExpectExec("INSERT INTO audit_log").WithArgs(
		"admin", "3", "order.payment_recorded", "order", "17", "req-1", "checkout-abc-123",
		"198.51.100.4", `{"paymentStatus":"pending"}`, `{"paymentStatus":"paid"}`, nil,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = NewAuditRepository(db).Record(ctx, nil, models.AuditEntry{
		Action: "order.payment_recorded", EntityType: "order", EntityID: "17",
		Before: map[string]string{"paymentStatus": "pending"},
		After:  map[string]string{"paymentStatus": "paid"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAuditRecordDefaultsToSystemActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO audit_log").WithArgs(
		"system", nil, "order.payment_expired", "order", "17", nil, nil, nil, nil, nil, nil,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = NewAuditRepository(db).Record(context.Background(), nil, models.AuditEntry{
		Action: "order.payment_expired", EntityType: "order", EntityID: "17",
	})
	if err != nil {
		t.Fatal(err)
	}
}
