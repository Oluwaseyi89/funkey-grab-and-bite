package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db, mock
}

func TestWithTxCommitsWhenFnSucceeds(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := WithTx(context.Background(), db, func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE orders SET status = 'confirmed'")
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTxRollsBackWhenFnFails(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectRollback()

	want := errors.New("out of stock")
	err := WithTx(context.Background(), db, func(tx *sql.Tx) error { return want })
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTxRollsBackAndRepanics(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectRollback()

	defer func() {
		if recover() == nil {
			t.Fatal("expected the panic to propagate")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	}()
	_ = WithTx(context.Background(), db, func(tx *sql.Tx) error { panic("boom") })
}

func TestWithTxRetriesSerializationFailures(t *testing.T) {
	db, mock := newMockDB(t)
	serialization := &pq.Error{Code: "40001"}
	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectCommit()

	attempts := 0
	err := WithTx(context.Background(), db, func(tx *sql.Tx) error {
		attempts++
		if attempts == 1 {
			return serialization
		}
		return nil
	})
	if err != nil || attempts != 2 {
		t.Fatalf("expected success on the second attempt, got %v after %d attempts", err, attempts)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestWithTxGivesUpAfterMaxAttempts(t *testing.T) {
	db, mock := newMockDB(t)
	for range maxTxAttempts {
		mock.ExpectBegin()
		mock.ExpectRollback()
	}

	attempts := 0
	err := WithTx(context.Background(), db, func(tx *sql.Tx) error {
		attempts++
		return &pq.Error{Code: "40P01"}
	})
	if err == nil || attempts != maxTxAttempts {
		t.Fatalf("expected failure after %d attempts, got %v after %d", maxTxAttempts, err, attempts)
	}
}

func TestWithTxDoesNotRetryOtherErrors(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectRollback()

	attempts := 0
	_ = WithTx(context.Background(), db, func(tx *sql.Tx) error {
		attempts++
		return &pq.Error{Code: "23505"} // unique violation
	})
	if attempts != 1 {
		t.Fatalf("expected a single attempt, got %d", attempts)
	}
}
