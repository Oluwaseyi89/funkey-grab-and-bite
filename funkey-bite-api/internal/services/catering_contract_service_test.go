package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// memoryContract keeps contract records in memory. Its transactions are not
// rolled back; the tests only check what successful calls leave behind.
type memoryContract struct {
	requests  *memoryCateringStore
	quotes    []models.CateringQuote
	hashes    map[int]string
	payments  []models.CateringPayment
	invoices  []models.CateringInvoice
	receipt   *models.CateringReceipt
	history   []models.CateringStatusChange
	charges   []models.CateringPaymentCharge
	emails    []string
	initiated int
}

func (m *memoryContract) InTx(fn func(tx repository.CateringContractWriter) error) error {
	return fn(m)
}

func (m *memoryContract) LockRequest(id int) (string, bool, error) {
	r, ok := m.requests.requests[id]
	if !ok {
		return "", false, nil
	}
	return string(r.Status), true, nil
}
func (m *memoryContract) SetStatus(requestID int, from, to, actorType string, actorID *int, note string) error {
	m.requests.requests[requestID].Status = models.CateringStatus(to)
	m.history = append(m.history, models.CateringStatusChange{RequestID: requestID, FromStatus: &from, ToStatus: to, ActorType: actorType, Note: &note})
	return nil
}
func (m *memoryContract) CloseOpenQuotes(requestID int, status string) error {
	for i := range m.quotes {
		if m.quotes[i].RequestID == requestID && m.quotes[i].Status == models.QuoteStatusSent {
			m.quotes[i].Status = status
		}
	}
	return nil
}
func (m *memoryContract) InsertQuote(quote *models.CateringQuote, tokenHash string, _ *int) error {
	quote.ID = len(m.quotes) + 1
	quote.Version = 1
	for _, q := range m.quotes {
		if q.RequestID == quote.RequestID {
			quote.Version++
		}
	}
	quote.Status = models.QuoteStatusSent
	m.quotes = append(m.quotes, *quote)
	m.hashes[quote.ID] = tokenHash
	amount := quote.Amount
	m.requests.requests[quote.RequestID].QuotedAmount = &amount
	return nil
}
func (m *memoryContract) AcceptQuote(quote *models.CateringQuote, via, invoiceNumber string, at time.Time) error {
	for i := range m.quotes {
		if m.quotes[i].ID == quote.ID {
			m.quotes[i].Status = models.QuoteStatusAccepted
			m.quotes[i].AcceptedVia = &via
		}
	}
	r := m.requests.requests[quote.RequestID]
	id := quote.ID
	r.AcceptedQuoteID = &id
	r.InvoiceNumber = &invoiceNumber
	return nil
}
func (m *memoryContract) GetQuote(id int) (*models.CateringQuote, error) {
	for _, q := range m.quotes {
		if q.ID == id {
			return &q, nil
		}
	}
	return nil, nil
}
func (m *memoryContract) InsertPayment(payment *models.CateringPayment) error {
	for _, p := range m.payments {
		if p.Reference != nil && payment.Reference != nil && *p.Reference == *payment.Reference && p.Method == models.CateringPaymentPaystackTransfer {
			return repository.ErrDuplicatePayment
		}
	}
	payment.ID = len(m.payments) + 1
	m.payments = append(m.payments, *payment)
	return nil
}
func (m *memoryContract) AmountPaid(requestID int) (float64, error) {
	total := 0.0
	for _, p := range m.payments {
		if p.RequestID == requestID {
			total += p.Amount
		}
	}
	return total, nil
}
func (m *memoryContract) InsertInvoice(invoice *models.CateringInvoice) error {
	invoice.Version = len(m.invoices) + 1
	m.invoices = append([]models.CateringInvoice{*invoice}, m.invoices...)
	return nil
}
func (m *memoryContract) InsertReceipt(receipt *models.CateringReceipt) error {
	m.receipt = receipt
	return nil
}
func (m *memoryContract) MarkContracted(requestID int, at time.Time) error {
	m.requests.requests[requestID].ContractedAt = &at
	return nil
}
func (m *memoryContract) MarkPaidInFull(requestID int, at time.Time) error {
	m.requests.requests[requestID].PaidInFullAt = &at
	return nil
}
func (m *memoryContract) InsertCharge(charge *models.CateringPaymentCharge) error {
	charge.ID = len(m.charges) + 1
	charge.Status = "pending"
	m.charges = append(m.charges, *charge)
	return nil
}
func (m *memoryContract) ChargeByReference(reference string) (*models.CateringPaymentCharge, error) {
	for _, c := range m.charges {
		if c.Reference == reference {
			return &c, nil
		}
	}
	return nil, nil
}
func (m *memoryContract) MarkChargePaid(id int) error {
	m.charges[id-1].Status = "paid"
	return nil
}

func (m *memoryContract) ListQuotes(requestID int) ([]models.CateringQuote, error) {
	var out []models.CateringQuote
	for i := len(m.quotes) - 1; i >= 0; i-- {
		if m.quotes[i].RequestID == requestID {
			out = append(out, m.quotes[i])
		}
	}
	return out, nil
}
func (m *memoryContract) QuoteByTokenHash(hash string) (*models.CateringQuote, error) {
	for id, h := range m.hashes {
		if h == hash {
			return m.GetQuote(id)
		}
	}
	return nil, nil
}
func (m *memoryContract) ListPayments(int) ([]models.CateringPayment, error) { return m.payments, nil }
func (m *memoryContract) ListInvoices(int) ([]models.CateringInvoice, error) { return m.invoices, nil }
func (m *memoryContract) GetReceipt(int) (*models.CateringReceipt, error)    { return m.receipt, nil }
func (m *memoryContract) ListHistory(int) ([]models.CateringStatusChange, error) {
	return m.history, nil
}
func (m *memoryContract) PendingCharge(requestID int, _ time.Time) (*models.CateringPaymentCharge, error) {
	for i := len(m.charges) - 1; i >= 0; i-- {
		if m.charges[i].RequestID == requestID && m.charges[i].Status == "pending" {
			c := m.charges[i]
			return &c, nil
		}
	}
	return nil, nil
}

type contractSettings struct{}

func (contractSettings) GetSettings() (*models.BusinessSettings, error) {
	return &models.BusinessSettings{BusinessName: "Funkey", PhoneNumber: "0800", CateringTerms: "Default terms"}, nil
}

type recordingEmail struct {
	utils.EmailService
	contract *memoryContract
}

func (e recordingEmail) SendInvoice(_, subject, _, _ string) error {
	e.contract.emails = append(e.contract.emails, subject)
	return nil
}

type fakeGateway struct {
	utils.PaymentGateway
	contract *memoryContract
}

func (g fakeGateway) InitiateTransfer(_ string, _ float64, reference string) (*utils.TransferChargeResult, error) {
	g.contract.initiated++
	return &utils.TransferChargeResult{Reference: reference, AccountNumber: "0123456789", BankName: "Test Bank"}, nil
}

func newTestContract(t *testing.T) (*cateringContractService, *memoryContract, *models.CateringRequest) {
	t.Helper()
	store := newMemoryCateringStore()
	email := "ada@example.com"
	request, _ := store.Create(&models.CateringRequest{
		Reference: "CAT-ABC123", ContactName: "Ada", ContactPhone: "08012345678", ContactEmail: &email,
		EventDate: "2026-11-20", GuestCount: 80, Status: models.CateringStatusPending,
	})
	mem := &memoryContract{requests: store, hashes: map[int]string{}}
	svc := NewCateringContractService(store, mem, contractSettings{}, recordingEmail{contract: mem}, fakeGateway{contract: mem}).(*cateringContractService)
	svc.now = func() time.Time { return cateringNow }
	svc.send = func(fn func()) { fn() }
	svc.webURL = "https://example.test"
	return svc, mem, request
}

func customerAuth() models.CateringCustomerAuth {
	return models.CateringCustomerAuth{Reference: "CAT-ABC123", Phone: "+234 801 234 5678"}
}

func TestSendQuoteDefaultsAndValidation(t *testing.T) {
	svc, mem, request := newTestContract(t)

	_, err := svc.SendQuote(request.ID, &models.CateringQuoteInput{
		Amount: 500000, LineItems: []models.CateringQuoteLine{{Description: "Food", Amount: 400000}},
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "add up to") {
		t.Fatalf("mismatched lines must be rejected, got %v", err)
	}

	contract, err := svc.SendQuote(request.ID, &models.CateringQuoteInput{Amount: 500000}, nil)
	if err != nil {
		t.Fatalf("SendQuote() error = %v", err)
	}
	q := contract.Quote
	if q == nil || q.DepositPercent != 70 || q.DepositAmount != 350000 || q.Terms != "Default terms" || q.ValidUntil != "2026-11-03" {
		t.Fatalf("unexpected quote defaults %+v", q)
	}
	if contract.Request.Status != models.CateringStatusQuoted || len(mem.emails) != 1 {
		t.Fatalf("expected quoted status and a quote email, got %s / %v", contract.Request.Status, mem.emails)
	}

	// Validity never runs past the event.
	days := 60
	contract, err = svc.SendQuote(request.ID, &models.CateringQuoteInput{Amount: 450000, ValidDays: &days}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if contract.Quote.Version != 2 || contract.Quote.ValidUntil != "2026-11-20" {
		t.Fatalf("expected v2 valid until the event, got %+v", contract.Quote)
	}
	if mem.quotes[0].Status != models.QuoteStatusSuperseded {
		t.Fatalf("older quote should be superseded, got %s", mem.quotes[0].Status)
	}
}

func TestAcceptQuoteAndPayInParts(t *testing.T) {
	svc, mem, request := newTestContract(t)
	if _, err := svc.SendQuote(request.ID, &models.CateringQuoteInput{Amount: 100000}, nil); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: customerAuth()}); err == nil {
		t.Fatal("accepting without agreeing to the terms must be rejected")
	}
	wrong := customerAuth()
	wrong.Phone = "08099999999"
	if _, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: wrong, AgreeToTerms: true}); !errors.Is(err, ErrCateringRequestNotFound) {
		t.Fatalf("wrong phone must not find the request, got %v", err)
	}

	contract, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: customerAuth(), AgreeToTerms: true})
	if err != nil {
		t.Fatalf("AcceptQuote() error = %v", err)
	}
	if contract.Request.Status != models.CateringStatusAccepted || contract.Invoice == nil || contract.Invoice.InvoiceNumber != "INV-CAT-ABC123" {
		t.Fatalf("expected accepted with invoice, got %+v", contract)
	}
	if contract.History != nil || contract.Quotes != nil || contract.Request.AdminNotes != nil {
		t.Fatal("customer view must hide history, quote list and notes")
	}
	if _, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: customerAuth(), AgreeToTerms: true}); err == nil {
		t.Fatal("accepting twice must be rejected")
	}

	if _, err := svc.RecordPayment(request.ID, &models.CateringPaymentInput{Amount: 200000, Method: "cash"}, nil); err == nil {
		t.Fatal("overpayment must be rejected")
	}
	if _, err := svc.RecordPayment(request.ID, &models.CateringPaymentInput{Amount: 1000, Method: "bank_transfer"}, nil); err == nil {
		t.Fatal("bank transfer without a reference must be rejected")
	}

	contract, err = svc.RecordPayment(request.ID, &models.CateringPaymentInput{Amount: 30000, Method: "cash"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if contract.Request.Status != models.CateringStatusAccepted || contract.Account.DepositRemaining != 40000 || contract.Invoice.Balance != 70000 {
		t.Fatalf("part of the deposit should re-issue the invoice, got %+v %+v", contract.Account, contract.Invoice)
	}

	contract, err = svc.RecordPayment(request.ID, &models.CateringPaymentInput{Amount: 40000, Method: "cash"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if contract.Request.Status != models.CateringStatusConfirmed || contract.Request.ContractedAt == nil || contract.Receipt != nil {
		t.Fatalf("deposit should confirm without a receipt, got %+v", contract.Request)
	}

	ref := "TRF-1"
	contract, err = svc.RecordPayment(request.ID, &models.CateringPaymentInput{Amount: 30000, Method: "bank_transfer", Reference: &ref}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if contract.Request.Status != models.CateringStatusPaid || contract.Receipt == nil || contract.Receipt.ReceiptNumber != "RCT-CAT-ABC123" || contract.Account.Balance != 0 {
		t.Fatalf("full payment should issue a receipt, got %+v", contract)
	}

	// quote, invoice on acceptance, two re-issued invoices, receipt
	want := []string{"Your catering quote", "Invoice INV-CAT-ABC123", "Invoice INV-CAT-ABC123", "Invoice INV-CAT-ABC123", "Receipt RCT-CAT-ABC123"}
	if len(mem.emails) != len(want) {
		t.Fatalf("emails = %v", mem.emails)
	}
	for i := range want {
		if !strings.HasPrefix(mem.emails[i], want[i]) {
			t.Fatalf("email %d = %q, want prefix %q", i, mem.emails[i], want[i])
		}
	}

	if _, err := svc.RecordPayment(request.ID, &models.CateringPaymentInput{Amount: 1, Method: "cash"}, nil); err == nil {
		t.Fatal("no payment is due once paid in full")
	}
}

func TestAcceptByEmailLinkAndExpiry(t *testing.T) {
	svc, mem, request := newTestContract(t)
	days := 1
	if _, err := svc.SendQuote(request.ID, &models.CateringQuoteInput{Amount: 100000, ValidDays: &days}, nil); err != nil {
		t.Fatal(err)
	}
	// The test reads the token back by replacing its stored hash.
	mem.hashes[1] = hashToken("link-token")

	if _, err := svc.CustomerView(models.CateringCustomerAuth{Token: "wrong"}); !errors.Is(err, ErrCateringLinkInvalid) {
		t.Fatalf("unknown token must be rejected, got %v", err)
	}

	svc.now = func() time.Time { return cateringNow.AddDate(0, 0, 2) }
	if _, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: models.CateringCustomerAuth{Token: "link-token"}, AgreeToTerms: true}); err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expired quote must be rejected, got %v", err)
	}

	svc.now = func() time.Time { return cateringNow }
	contract, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: models.CateringCustomerAuth{Token: "link-token"}, AgreeToTerms: true})
	if err != nil {
		t.Fatal(err)
	}
	if contract.Quote.AcceptedVia == nil || *contract.Quote.AcceptedVia != models.AcceptedViaEmailLink {
		t.Fatalf("expected acceptance by email link, got %+v", contract.Quote)
	}
}

func TestOnlinePaymentWebhook(t *testing.T) {
	svc, mem, request := newTestContract(t)
	if _, err := svc.SendQuote(request.ID, &models.CateringQuoteInput{Amount: 100000}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AcceptQuote(models.CateringAcceptInput{CateringCustomerAuth: customerAuth(), AgreeToTerms: true}); err != nil {
		t.Fatal(err)
	}

	charge, err := svc.StartOnlinePayment(models.CateringPayOnlineInput{CateringCustomerAuth: customerAuth(), Portion: "deposit"})
	if err != nil {
		t.Fatal(err)
	}
	if charge.Amount != 70000 || !strings.HasPrefix(charge.Reference, CateringPaystackPrefix) {
		t.Fatalf("unexpected charge %+v", charge)
	}
	again, _ := svc.StartOnlinePayment(models.CateringPayOnlineInput{CateringCustomerAuth: customerAuth(), Portion: "deposit"})
	if again.Reference != charge.Reference || mem.initiated != 1 {
		t.Fatal("a pending charge for the same amount should be reused")
	}

	if handled, _ := svc.HandlePaystackCharge("FGB-ORDER-1", 100); handled {
		t.Fatal("order references must be left to the order flow")
	}
	if handled, err := svc.HandlePaystackCharge(charge.Reference, 7000000); !handled || err != nil {
		t.Fatalf("HandlePaystackCharge() = %v, %v", handled, err)
	}
	if handled, err := svc.HandlePaystackCharge(charge.Reference, 7000000); !handled || err != nil {
		t.Fatalf("duplicate delivery should be ignored, got %v, %v", handled, err)
	}

	contract, _ := svc.GetContract(request.ID)
	if contract.Request.Status != models.CateringStatusConfirmed || len(contract.Payments) != 1 || contract.Payments[0].Method != models.CateringPaymentPaystackTransfer {
		t.Fatalf("expected one Paystack payment confirming the booking, got %+v", contract)
	}
}

func TestChangeStatusRules(t *testing.T) {
	svc, mem, request := newTestContract(t)

	if _, err := svc.ChangeStatus(request.ID, "paid", nil, ""); err == nil {
		t.Fatal("paid must only be set by payments")
	}
	if _, err := svc.ChangeStatus(request.ID, "completed", nil, ""); err == nil {
		t.Fatal("only a paid booking can be completed")
	}
	if _, err := svc.SendQuote(request.ID, &models.CateringQuoteInput{Amount: 1000}, nil); err != nil {
		t.Fatal(err)
	}
	contract, err := svc.ChangeStatus(request.ID, "declined", nil, "Fully booked")
	if err != nil {
		t.Fatal(err)
	}
	if contract.Request.Status != models.CateringStatusDeclined || mem.quotes[0].Status != models.QuoteStatusWithdrawn {
		t.Fatalf("declining should withdraw the open quote, got %s / %s", contract.Request.Status, mem.quotes[0].Status)
	}
	if _, err := svc.ChangeStatus(request.ID, "pending", nil, ""); err != nil {
		t.Fatalf("a declined request without an accepted quote can be reopened: %v", err)
	}
}
