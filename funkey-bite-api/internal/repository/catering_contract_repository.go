package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/lib/pq"
)

// ErrDuplicatePayment is returned when an online payment reference was
// already recorded (a webhook delivered twice).
var ErrDuplicatePayment = errors.New("payment already recorded")

// CateringContractRepository stores quotes, payments, invoices, receipts and
// status history for catering requests.
type CateringContractRepository struct {
	db *sql.DB
}

func NewCateringContractRepository(db *sql.DB) *CateringContractRepository {
	return &CateringContractRepository{db: db}
}

// CateringContractWriter is what a contract transaction can do; the service
// depends on this so its rules can be tested without a database.
type CateringContractWriter interface {
	LockRequest(id int) (status string, ok bool, err error)
	SetStatus(requestID int, from, to, actorType string, actorID *int, note string) error
	CloseOpenQuotes(requestID int, status string) error
	InsertQuote(quote *models.CateringQuote, tokenHash string, createdBy *int) error
	AcceptQuote(quote *models.CateringQuote, via, invoiceNumber string, at time.Time) error
	GetQuote(id int) (*models.CateringQuote, error)
	InsertPayment(payment *models.CateringPayment) error
	AmountPaid(requestID int) (float64, error)
	InsertInvoice(invoice *models.CateringInvoice) error
	InsertReceipt(receipt *models.CateringReceipt) error
	MarkContracted(requestID int, at time.Time) error
	MarkPaidInFull(requestID int, at time.Time) error
	InsertCharge(charge *models.CateringPaymentCharge) error
	ChargeByReference(reference string) (*models.CateringPaymentCharge, error)
	MarkChargePaid(id int) error
}

// CateringContractTx runs contract writes inside one database transaction,
// with the request row locked, so a status and the records that justify it
// are always saved together.
type CateringContractTx struct {
	tx *sql.Tx
}

// InTx runs fn in a transaction, committing only if it returns nil.
func (r *CateringContractRepository) InTx(fn func(tx CateringContractWriter) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	if err := fn(&CateringContractTx{tx: tx}); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// LockRequest locks a request for the rest of the transaction and returns
// its status; ok is false when there is no such request.
func (t *CateringContractTx) LockRequest(id int) (status string, ok bool, err error) {
	err = t.tx.QueryRow(`SELECT status FROM catering_requests WHERE id = $1 FOR UPDATE`, id).Scan(&status)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to lock catering request: %w", err)
	}
	return status, true, nil
}

// SetStatus changes a request's status and records the change.
func (t *CateringContractTx) SetStatus(requestID int, from, to, actorType string, actorID *int, note string) error {
	if _, err := t.tx.Exec(
		`UPDATE catering_requests SET status = $1, updated_at = $2 WHERE id = $3`,
		to, time.Now(), requestID,
	); err != nil {
		return fmt.Errorf("failed to update catering status: %w", err)
	}
	var fromValue, noteValue *string
	if from != "" {
		fromValue = &from
	}
	if note != "" {
		noteValue = &note
	}
	if _, err := t.tx.Exec(`
		INSERT INTO catering_status_history (request_id, from_status, to_status, actor_type, actor_id, note)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		requestID, fromValue, to, actorType, actorID, noteValue,
	); err != nil {
		return fmt.Errorf("failed to record catering status change: %w", err)
	}
	return nil
}

// CloseOpenQuotes marks quotes still awaiting the customer as status
// (superseded by a newer quote, or withdrawn when the request closes).
func (t *CateringContractTx) CloseOpenQuotes(requestID int, status string) error {
	if _, err := t.tx.Exec(
		`UPDATE catering_quotes SET status = $1, acceptance_token_hash = NULL WHERE request_id = $2 AND status = 'sent'`,
		status, requestID,
	); err != nil {
		return fmt.Errorf("failed to close open quotes: %w", err)
	}
	return nil
}

// InsertQuote saves a new quote version and mirrors its amount on the request.
func (t *CateringContractTx) InsertQuote(quote *models.CateringQuote, tokenHash string, createdBy *int) error {
	if err := t.tx.QueryRow(
		`SELECT COALESCE(MAX(version), 0) + 1 FROM catering_quotes WHERE request_id = $1`, quote.RequestID,
	).Scan(&quote.Version); err != nil {
		return fmt.Errorf("failed to number quote: %w", err)
	}

	lines, err := json.Marshal(quote.LineItems)
	if err != nil {
		return fmt.Errorf("failed to encode quote lines: %w", err)
	}
	err = t.tx.QueryRow(`
		INSERT INTO catering_quotes (
			request_id, version, amount, line_items, deposit_percent, terms,
			message, valid_until, status, acceptance_token_hash, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'sent', $9, $10)
		RETURNING id, created_at`,
		quote.RequestID, quote.Version, quote.Amount, lines, quote.DepositPercent, quote.Terms,
		quote.Message, quote.ValidUntil, tokenHash, createdBy,
	).Scan(&quote.ID, &quote.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save quote: %w", err)
	}
	quote.Status = models.QuoteStatusSent

	if _, err := t.tx.Exec(
		`UPDATE catering_requests SET quoted_amount = $1, updated_at = $2 WHERE id = $3`,
		quote.Amount, time.Now(), quote.RequestID,
	); err != nil {
		return fmt.Errorf("failed to record quoted amount: %w", err)
	}
	return nil
}

// AcceptQuote marks a quote accepted and links it and the invoice number to
// its request.
func (t *CateringContractTx) AcceptQuote(quote *models.CateringQuote, via, invoiceNumber string, at time.Time) error {
	if _, err := t.tx.Exec(`
		UPDATE catering_quotes
		SET status = 'accepted', accepted_at = $1, accepted_via = $2
		WHERE id = $3`,
		at, via, quote.ID,
	); err != nil {
		return fmt.Errorf("failed to accept quote: %w", err)
	}
	if _, err := t.tx.Exec(`
		UPDATE catering_requests SET accepted_quote_id = $1, invoice_number = $2, updated_at = $3
		WHERE id = $4`,
		quote.ID, invoiceNumber, at, quote.RequestID,
	); err != nil {
		return fmt.Errorf("failed to link accepted quote: %w", err)
	}
	return nil
}

// GetQuote reads a quote inside the transaction.
func (t *CateringContractTx) GetQuote(id int) (*models.CateringQuote, error) {
	quote, err := scanQuote(t.tx.QueryRow(`SELECT `+quoteColumns+` FROM catering_quotes WHERE id = $1`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return quote, err
}

// InsertPayment adds a payment to the ledger.
func (t *CateringContractTx) InsertPayment(payment *models.CateringPayment) error {
	err := t.tx.QueryRow(`
		INSERT INTO catering_payments (request_id, amount, method, reference, note, recorded_by, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`,
		payment.RequestID, payment.Amount, payment.Method, payment.Reference, payment.Note,
		payment.RecordedBy, payment.PaidAt,
	).Scan(&payment.ID, &payment.CreatedAt)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrDuplicatePayment
	}
	if err != nil {
		return fmt.Errorf("failed to record payment: %w", err)
	}
	return nil
}

// AmountPaid sums the ledger for a request.
func (t *CateringContractTx) AmountPaid(requestID int) (float64, error) {
	var total float64
	if err := t.tx.QueryRow(
		`SELECT COALESCE(SUM(amount), 0) FROM catering_payments WHERE request_id = $1`, requestID,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("failed to total payments: %w", err)
	}
	return total, nil
}

// InsertInvoice issues the next version of a booking's invoice.
func (t *CateringContractTx) InsertInvoice(invoice *models.CateringInvoice) error {
	if err := t.tx.QueryRow(
		`SELECT COALESCE(MAX(version), 0) + 1 FROM catering_invoices WHERE invoice_number = $1`, invoice.InvoiceNumber,
	).Scan(&invoice.Version); err != nil {
		return fmt.Errorf("failed to number invoice: %w", err)
	}
	if err := t.tx.QueryRow(`
		INSERT INTO catering_invoices (
			request_id, quote_id, invoice_number, version, reason, total,
			deposit_due, amount_paid, balance
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, issued_at`,
		invoice.RequestID, invoice.QuoteID, invoice.InvoiceNumber, invoice.Version, invoice.Reason,
		invoice.Total, invoice.DepositDue, invoice.AmountPaid, invoice.Balance,
	).Scan(&invoice.ID, &invoice.IssuedAt); err != nil {
		return fmt.Errorf("failed to issue invoice: %w", err)
	}
	return nil
}

// InsertReceipt issues the booking's receipt.
func (t *CateringContractTx) InsertReceipt(receipt *models.CateringReceipt) error {
	if err := t.tx.QueryRow(`
		INSERT INTO catering_receipts (request_id, receipt_number, total_paid)
		VALUES ($1, $2, $3)
		RETURNING id, issued_at`,
		receipt.RequestID, receipt.ReceiptNumber, receipt.TotalPaid,
	).Scan(&receipt.ID, &receipt.IssuedAt); err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}
	return nil
}

// MarkContracted records when the booking became a contract.
func (t *CateringContractTx) MarkContracted(requestID int, at time.Time) error {
	_, err := t.tx.Exec(`UPDATE catering_requests SET contracted_at = $1 WHERE id = $2 AND contracted_at IS NULL`, at, requestID)
	return err
}

// MarkPaidInFull records when the booking was paid in full.
func (t *CateringContractTx) MarkPaidInFull(requestID int, at time.Time) error {
	_, err := t.tx.Exec(`UPDATE catering_requests SET paid_in_full_at = $1 WHERE id = $2 AND paid_in_full_at IS NULL`, at, requestID)
	return err
}

// InsertCharge saves a Paystack dedicated account awaiting a transfer.
func (t *CateringContractTx) InsertCharge(charge *models.CateringPaymentCharge) error {
	return t.tx.QueryRow(`
		INSERT INTO catering_payment_charges (
			request_id, reference, amount, account_number, account_name, bank_name, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, status, created_at`,
		charge.RequestID, charge.Reference, charge.Amount, charge.AccountNumber,
		charge.AccountName, charge.BankName, charge.ExpiresAt,
	).Scan(&charge.ID, &charge.Status, &charge.CreatedAt)
}

// ChargeByReference finds a charge and locks it; nil when there is none.
func (t *CateringContractTx) ChargeByReference(reference string) (*models.CateringPaymentCharge, error) {
	charge, err := scanCharge(t.tx.QueryRow(`SELECT `+chargeColumns+` FROM catering_payment_charges WHERE reference = $1 FOR UPDATE`, reference))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return charge, err
}

// MarkChargePaid closes a charge once its transfer is recorded.
func (t *CateringContractTx) MarkChargePaid(id int) error {
	_, err := t.tx.Exec(`UPDATE catering_payment_charges SET status = 'paid' WHERE id = $1`, id)
	return err
}

const quoteColumns = `id, request_id, version, amount, line_items, deposit_percent, terms,
	message, valid_until, status, accepted_at, accepted_via, created_at`

func scanQuote(row rowScanner) (*models.CateringQuote, error) {
	var quote models.CateringQuote
	var lines []byte
	var validUntil time.Time
	var acceptedAt sql.NullTime
	err := row.Scan(&quote.ID, &quote.RequestID, &quote.Version, &quote.Amount, &lines, &quote.DepositPercent,
		&quote.Terms, &quote.Message, &validUntil, &quote.Status, &acceptedAt, &quote.AcceptedVia, &quote.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan quote: %w", err)
	}
	quote.ValidUntil = validUntil.Format("2006-01-02")
	if acceptedAt.Valid {
		quote.AcceptedAt = &acceptedAt.Time
	}
	quote.LineItems = []models.CateringQuoteLine{}
	if len(lines) > 0 {
		if err := json.Unmarshal(lines, &quote.LineItems); err != nil {
			return nil, fmt.Errorf("failed to parse quote lines: %w", err)
		}
	}
	return &quote, nil
}

// ListQuotes returns every quote version for a request, newest first.
func (r *CateringContractRepository) ListQuotes(requestID int) ([]models.CateringQuote, error) {
	rows, err := r.db.Query(`SELECT `+quoteColumns+` FROM catering_quotes WHERE request_id = $1 ORDER BY version DESC`, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	quotes := []models.CateringQuote{}
	for rows.Next() {
		quote, err := scanQuote(rows)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, *quote)
	}
	return quotes, rows.Err()
}

// QuoteByTokenHash finds the open quote an emailed link belongs to.
func (r *CateringContractRepository) QuoteByTokenHash(tokenHash string) (*models.CateringQuote, error) {
	quote, err := scanQuote(r.db.QueryRow(
		`SELECT `+quoteColumns+` FROM catering_quotes WHERE acceptance_token_hash = $1`, tokenHash))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return quote, err
}

func (r *CateringContractRepository) ListPayments(requestID int) ([]models.CateringPayment, error) {
	rows, err := r.db.Query(`
		SELECT id, request_id, amount, method, reference, note, recorded_by, paid_at, created_at
		FROM catering_payments WHERE request_id = $1 ORDER BY paid_at, id`, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	payments := []models.CateringPayment{}
	for rows.Next() {
		var p models.CateringPayment
		var recordedBy sql.NullInt64
		if err := rows.Scan(&p.ID, &p.RequestID, &p.Amount, &p.Method, &p.Reference, &p.Note,
			&recordedBy, &p.PaidAt, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		p.RecordedBy = nullIntPtr(recordedBy)
		payments = append(payments, p)
	}
	return payments, rows.Err()
}

// ListInvoices returns every issued invoice version, newest first.
func (r *CateringContractRepository) ListInvoices(requestID int) ([]models.CateringInvoice, error) {
	rows, err := r.db.Query(`
		SELECT id, request_id, quote_id, invoice_number, version, reason, total,
		       deposit_due, amount_paid, balance, issued_at
		FROM catering_invoices WHERE request_id = $1 ORDER BY version DESC`, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	defer rows.Close()

	invoices := []models.CateringInvoice{}
	for rows.Next() {
		var inv models.CateringInvoice
		if err := rows.Scan(&inv.ID, &inv.RequestID, &inv.QuoteID, &inv.InvoiceNumber, &inv.Version,
			&inv.Reason, &inv.Total, &inv.DepositDue, &inv.AmountPaid, &inv.Balance, &inv.IssuedAt); err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		invoices = append(invoices, inv)
	}
	return invoices, rows.Err()
}

func (r *CateringContractRepository) GetReceipt(requestID int) (*models.CateringReceipt, error) {
	var receipt models.CateringReceipt
	err := r.db.QueryRow(`
		SELECT id, request_id, receipt_number, total_paid, issued_at
		FROM catering_receipts WHERE request_id = $1`, requestID,
	).Scan(&receipt.ID, &receipt.RequestID, &receipt.ReceiptNumber, &receipt.TotalPaid, &receipt.IssuedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}
	return &receipt, nil
}

func (r *CateringContractRepository) ListHistory(requestID int) ([]models.CateringStatusChange, error) {
	rows, err := r.db.Query(`
		SELECT id, request_id, from_status, to_status, actor_type, actor_id, note, created_at
		FROM catering_status_history WHERE request_id = $1 ORDER BY created_at, id`, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to list status history: %w", err)
	}
	defer rows.Close()

	history := []models.CateringStatusChange{}
	for rows.Next() {
		var change models.CateringStatusChange
		var actorID sql.NullInt64
		if err := rows.Scan(&change.ID, &change.RequestID, &change.FromStatus, &change.ToStatus,
			&change.ActorType, &actorID, &change.Note, &change.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan status change: %w", err)
		}
		change.ActorID = nullIntPtr(actorID)
		history = append(history, change)
	}
	return history, rows.Err()
}

// PendingCharge returns the newest unpaid, unexpired Paystack account for a
// request, so the customer is shown the same account again.
func (r *CateringContractRepository) PendingCharge(requestID int, now time.Time) (*models.CateringPaymentCharge, error) {
	charge, err := scanCharge(r.db.QueryRow(`
		SELECT `+chargeColumns+` FROM catering_payment_charges
		WHERE request_id = $1 AND status = 'pending' AND (expires_at IS NULL OR expires_at > $2)
		ORDER BY created_at DESC LIMIT 1`, requestID, now))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return charge, err
}

const chargeColumns = `id, request_id, reference, amount, account_number, account_name,
	bank_name, expires_at, status, created_at`

func scanCharge(row rowScanner) (*models.CateringPaymentCharge, error) {
	var c models.CateringPaymentCharge
	var expires sql.NullTime
	err := row.Scan(&c.ID, &c.RequestID, &c.Reference, &c.Amount, &c.AccountNumber, &c.AccountName,
		&c.BankName, &expires, &c.Status, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan payment charge: %w", err)
	}
	if expires.Valid {
		c.ExpiresAt = &expires.Time
	}
	return &c, nil
}
