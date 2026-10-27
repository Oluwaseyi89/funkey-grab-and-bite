package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// DefaultQuoteValidDays is how long a quote can be accepted when the admin
// does not set a validity.
const DefaultQuoteValidDays = 7

// CateringPaystackPrefix starts every Paystack reference for a catering
// payment, so the shared webhook can route it here.
const CateringPaystackPrefix = "CATPAY-"

// Actor types recorded in the status history.
const (
	ActorAdmin    = "admin"
	ActorCustomer = "customer"
	ActorSystem   = "system"
)

var ErrCateringLinkInvalid = errors.New("this link is no longer valid")

// CateringContractService runs a catering request from quote to payment:
//
//	pending → quoted (admin sends a quote)
//	        → accepted (customer accepts it; invoice issued)
//	        → confirmed (payments reach the deposit; the contract)
//	        → paid (payments reach the total; receipt issued)
//	        → completed (admin, after the event)
//
// Part payments short of the total re-issue the invoice with the balance.
type CateringContractService interface {
	GetContract(requestID int) (*models.CateringContract, error)
	SendQuote(requestID int, input *models.CateringQuoteInput, adminID *int) (*models.CateringContract, error)
	RecordPayment(requestID int, input *models.CateringPaymentInput, adminID *int) (*models.CateringContract, error)
	ChangeStatus(requestID int, to string, adminID *int, note string) (*models.CateringContract, error)

	CustomerView(auth models.CateringCustomerAuth) (*models.CateringContract, error)
	AcceptQuote(input models.CateringAcceptInput) (*models.CateringContract, error)
	StartOnlinePayment(input models.CateringPayOnlineInput) (*models.CateringPaymentCharge, error)

	// HandlePaystackCharge records a confirmed Paystack transfer; handled is
	// false when the reference is not a catering payment.
	HandlePaystackCharge(reference string, amountKobo int64) (handled bool, err error)
}

type cateringContractReader interface {
	InTx(fn func(tx repository.CateringContractWriter) error) error
	ListQuotes(requestID int) ([]models.CateringQuote, error)
	QuoteByTokenHash(tokenHash string) (*models.CateringQuote, error)
	ListPayments(requestID int) ([]models.CateringPayment, error)
	ListInvoices(requestID int) ([]models.CateringInvoice, error)
	GetReceipt(requestID int) (*models.CateringReceipt, error)
	ListHistory(requestID int) ([]models.CateringStatusChange, error)
	PendingCharge(requestID int, now time.Time) (*models.CateringPaymentCharge, error)
}

type cateringRequestReader interface {
	GetByID(id int) (*models.CateringRequest, error)
	GetByReference(reference string) (*models.CateringRequest, error)
}

type cateringSettingsReader interface {
	GetSettings() (*models.BusinessSettings, error)
}

type cateringContractService struct {
	requests cateringRequestReader
	contract cateringContractReader
	settings cateringSettingsReader
	email    utils.EmailService
	gateway  utils.PaymentGateway
	webURL   string
	location *time.Location
	now      func() time.Time
	// send runs email delivery; tests run it inline.
	send func(func())
}

func NewCateringContractService(
	requests cateringRequestReader,
	contract cateringContractReader,
	settings cateringSettingsReader,
	email utils.EmailService,
	gateway utils.PaymentGateway,
) CateringContractService {
	return &cateringContractService{
		requests: requests,
		contract: contract,
		settings: settings,
		email:    email,
		gateway:  gateway,
		// The website's address (as for password reset links), for the
		// booking links in emails.
		webURL:   strings.TrimRight(strings.TrimSpace(os.Getenv("APP_URL")), "/"),
		location: LoadBusinessLocation(),
		now:      time.Now,
		send:     func(fn func()) { go fn() },
	}
}

func roundNaira(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func depositFor(quote *models.CateringQuote) float64 {
	return roundNaira(quote.Amount * float64(quote.DepositPercent) / 100)
}

func (s *cateringContractService) today() time.Time {
	now := s.now().In(s.location)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.location)
}

func (s *cateringContractService) quoteExpired(quote *models.CateringQuote) bool {
	if quote.Status != models.QuoteStatusSent {
		return false
	}
	validUntil, err := time.ParseInLocation("2006-01-02", quote.ValidUntil, s.location)
	return err == nil && s.today().After(validUntil)
}

// SendQuote prices a pending or quoted request. A newer quote replaces any
// quote the customer has not accepted yet.
func (s *cateringContractService) SendQuote(requestID int, input *models.CateringQuoteInput, adminID *int) (*models.CateringContract, error) {
	request, err := s.request(requestID)
	if err != nil {
		return nil, err
	}

	quote := &models.CateringQuote{
		RequestID:      requestID,
		Amount:         roundNaira(input.Amount),
		LineItems:      []models.CateringQuoteLine{},
		DepositPercent: models.DefaultCateringDepositPercent,
		Message:        trimmedPtr(input.Message),
	}
	if input.DepositPercent != nil {
		quote.DepositPercent = *input.DepositPercent
	}

	if len(input.LineItems) > 0 {
		sum := 0.0
		for _, line := range input.LineItems {
			description := strings.TrimSpace(line.Description)
			if description == "" {
				return nil, invalid("Every quote line needs a description")
			}
			sum += line.Amount
			quote.LineItems = append(quote.LineItems, models.CateringQuoteLine{Description: description, Amount: roundNaira(line.Amount)})
		}
		if math.Abs(roundNaira(sum)-quote.Amount) > 0.009 {
			return nil, invalid("The quote lines add up to %s, not the quoted %s", FormatNaira(sum), FormatNaira(quote.Amount))
		}
	}

	terms := ""
	if input.Terms != nil {
		terms = strings.TrimSpace(*input.Terms)
	} else if settings, err := s.settings.GetSettings(); err == nil && settings != nil {
		terms = strings.TrimSpace(settings.CateringTerms)
	}
	if terms == "" {
		return nil, invalid("Enter the terms the customer agrees to")
	}
	quote.Terms = terms

	eventDate, err := time.ParseInLocation("2006-01-02", request.EventDate, s.location)
	if err != nil || eventDate.Before(s.today()) {
		return nil, invalid("The event date has passed, so it can no longer be quoted")
	}
	validDays := DefaultQuoteValidDays
	if input.ValidDays != nil {
		validDays = *input.ValidDays
	}
	validUntil := s.today().AddDate(0, 0, validDays)
	if validUntil.After(eventDate) {
		validUntil = eventDate
	}
	quote.ValidUntil = validUntil.Format("2006-01-02")

	token, tokenHash, err := newQuoteToken()
	if err != nil {
		return nil, err
	}

	err = s.contract.InTx(func(tx repository.CateringContractWriter) error {
		status, ok, err := tx.LockRequest(requestID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCateringRequestNotFound
		}
		if status != string(models.CateringStatusPending) && status != string(models.CateringStatusQuoted) {
			return invalid("A quote can only be sent before the customer accepts one")
		}
		if err := tx.CloseOpenQuotes(requestID, models.QuoteStatusSuperseded); err != nil {
			return err
		}
		if err := tx.InsertQuote(quote, tokenHash, adminID); err != nil {
			return err
		}
		return tx.SetStatus(requestID, status, string(models.CateringStatusQuoted), ActorAdmin, adminID,
			fmt.Sprintf("Quote v%d sent: %s", quote.Version, FormatNaira(quote.Amount)))
	})
	if err != nil {
		return nil, err
	}

	s.notify(func() { s.sendQuoteEmail(request, quote, token) })
	return s.GetContract(requestID)
}

// AcceptQuote records the customer's acceptance of the open quote and issues
// the invoice.
func (s *cateringContractService) AcceptQuote(input models.CateringAcceptInput) (*models.CateringContract, error) {
	if !input.AgreeToTerms {
		return nil, invalid("Please read and agree to the terms to accept the quote")
	}
	request, linkQuote, via, err := s.resolveCustomer(input.CateringCustomerAuth)
	if err != nil {
		return nil, err
	}

	quoteID := 0
	if linkQuote != nil {
		quoteID = linkQuote.ID
	} else if current, err := s.currentQuote(request.ID); err != nil {
		return nil, err
	} else if current != nil {
		quoteID = current.ID
	}
	if quoteID == 0 {
		return nil, invalid("There is no quote to accept yet")
	}

	invoiceNumber := "INV-" + request.Reference
	var invoice *models.CateringInvoice

	err = s.contract.InTx(func(tx repository.CateringContractWriter) error {
		status, ok, err := tx.LockRequest(request.ID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCateringRequestNotFound
		}
		quote, err := tx.GetQuote(quoteID)
		if err != nil {
			return err
		}
		if quote == nil || quote.RequestID != request.ID {
			return ErrCateringLinkInvalid
		}
		if quote.Status == models.QuoteStatusAccepted {
			return invalid("This quote has already been accepted")
		}
		if quote.Status != models.QuoteStatusSent || status != string(models.CateringStatusQuoted) {
			return invalid("This quote is no longer open; please contact us")
		}
		if s.quoteExpired(quote) {
			return invalid("This quote expired on %s; please contact us for a new one", quote.ValidUntil)
		}

		now := s.now()
		if err := tx.AcceptQuote(quote, via, invoiceNumber, now); err != nil {
			return err
		}
		if err := tx.SetStatus(request.ID, status, string(models.CateringStatusAccepted), ActorCustomer, nil,
			fmt.Sprintf("Quote v%d accepted (%s)", quote.Version, strings.ReplaceAll(via, "_", " "))); err != nil {
			return err
		}
		invoice = &models.CateringInvoice{
			RequestID: request.ID, QuoteID: quote.ID, InvoiceNumber: invoiceNumber, Reason: "accepted",
			Total: quote.Amount, DepositDue: depositFor(quote), AmountPaid: 0, Balance: quote.Amount,
		}
		return tx.InsertInvoice(invoice)
	})
	if err != nil {
		return nil, err
	}

	contract, err := s.GetContract(request.ID)
	if err != nil {
		return nil, err
	}
	s.notify(func() { s.sendInvoiceEmail(contract, invoice) })
	return customerContract(contract), nil
}

// RecordPayment adds a cash or direct bank transfer payment taken by the admin.
func (s *cateringContractService) RecordPayment(requestID int, input *models.CateringPaymentInput, adminID *int) (*models.CateringContract, error) {
	paidAt := s.now()
	if input.PaidOn != nil && strings.TrimSpace(*input.PaidOn) != "" {
		day, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(*input.PaidOn), s.location)
		if err != nil {
			return nil, invalid("Payment date must be in YYYY-MM-DD format")
		}
		if day.After(s.today()) {
			return nil, invalid("Payment date cannot be in the future")
		}
		if !day.Equal(s.today()) {
			paidAt = day.Add(12 * time.Hour)
		}
	}
	if input.Method == models.CateringPaymentBankTransfer && (input.Reference == nil || strings.TrimSpace(*input.Reference) == "") {
		return nil, invalid("Enter the bank transfer reference")
	}

	payment := &models.CateringPayment{
		RequestID:  requestID,
		Amount:     roundNaira(input.Amount),
		Method:     input.Method,
		Reference:  trimmedPtr(input.Reference),
		Note:       trimmedPtr(input.Note),
		RecordedBy: adminID,
		PaidAt:     paidAt,
	}

	var outcome *paymentOutcome
	err := s.contract.InTx(func(tx repository.CateringContractWriter) error {
		var err error
		outcome, err = s.applyPayment(tx, requestID, payment, ActorAdmin, adminID, true)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.afterPayment(requestID, outcome)
}

// HandlePaystackCharge records a Paystack transfer confirmed by the webhook.
// Deliveries can repeat, so an already-paid charge is ignored.
func (s *cateringContractService) HandlePaystackCharge(reference string, amountKobo int64) (bool, error) {
	if !strings.HasPrefix(reference, CateringPaystackPrefix) {
		return false, nil
	}

	var outcome *paymentOutcome
	var requestID int
	err := s.contract.InTx(func(tx repository.CateringContractWriter) error {
		charge, err := tx.ChargeByReference(reference)
		if err != nil {
			return err
		}
		if charge == nil {
			return fmt.Errorf("no catering payment found for reference %q", reference)
		}
		if charge.Status == "paid" {
			log.Printf("[Paystack webhook] catering payment %s already recorded, ignoring duplicate delivery", reference)
			return nil
		}
		requestID = charge.RequestID

		amount := roundNaira(float64(amountKobo) / 100)
		if math.Abs(amount-charge.Amount) > 0.009 {
			// Money was received either way, so it is recorded; the admin sees
			// the difference in the ledger.
			log.Printf("[Paystack webhook] catering payment %s: expected %s, received %s", reference, FormatNaira(charge.Amount), FormatNaira(amount))
		}
		payment := &models.CateringPayment{
			RequestID: charge.RequestID,
			Amount:    amount,
			Method:    models.CateringPaymentPaystackTransfer,
			Reference: &reference,
			PaidAt:    s.now(),
		}
		outcome, err = s.applyPayment(tx, charge.RequestID, payment, ActorSystem, nil, false)
		if err != nil {
			return err
		}
		return tx.MarkChargePaid(charge.ID)
	})
	if errors.Is(err, repository.ErrDuplicatePayment) {
		log.Printf("[Paystack webhook] catering payment %s already recorded, ignoring duplicate delivery", reference)
		return true, nil
	}
	if err != nil || outcome == nil {
		return true, err
	}
	_, err = s.afterPayment(requestID, outcome)
	return true, err
}

// paymentOutcome says which document a payment produced.
type paymentOutcome struct {
	invoice *models.CateringInvoice
	receipt *models.CateringReceipt
}

// applyPayment records a payment and moves the request on: to confirmed once
// the deposit is covered, and to paid with a receipt once the total is.
// Otherwise it re-issues the invoice with the new balance.
func (s *cateringContractService) applyPayment(
	tx repository.CateringContractWriter, requestID int, payment *models.CateringPayment,
	actorType string, actorID *int, strict bool,
) (*paymentOutcome, error) {
	status, ok, err := tx.LockRequest(requestID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCateringRequestNotFound
	}

	request, err := s.requests.GetByID(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil || request.AcceptedQuoteID == nil {
		return nil, invalid("Payments can only be recorded after the customer accepts a quote")
	}
	open := status == string(models.CateringStatusAccepted) || status == string(models.CateringStatusConfirmed)
	if strict && !open {
		return nil, invalid("This booking is %s, so no payment is due", status)
	}

	quote, err := tx.GetQuote(*request.AcceptedQuoteID)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, fmt.Errorf("accepted quote %d is missing", *request.AcceptedQuoteID)
	}

	paidBefore, err := tx.AmountPaid(requestID)
	if err != nil {
		return nil, err
	}
	balance := roundNaira(quote.Amount - paidBefore)
	if strict && payment.Amount > balance+0.009 {
		return nil, invalid("That is more than the %s still owed", FormatNaira(balance))
	}

	if err := tx.InsertPayment(payment); err != nil {
		return nil, err
	}
	paid := roundNaira(paidBefore + payment.Amount)
	note := fmt.Sprintf("%s received by %s", FormatNaira(payment.Amount), strings.ReplaceAll(payment.Method, "_", " "))
	if !open {
		log.Printf("Catering payment recorded for %s while it is %s", request.Reference, status)
		return &paymentOutcome{}, nil
	}

	now := s.now()
	if status == string(models.CateringStatusAccepted) && paid+0.009 >= depositFor(quote) {
		if err := tx.SetStatus(requestID, status, string(models.CateringStatusConfirmed), actorType, actorID, note+"; deposit reached, booking confirmed"); err != nil {
			return nil, err
		}
		if err := tx.MarkContracted(requestID, now); err != nil {
			return nil, err
		}
		status = string(models.CateringStatusConfirmed)
	}

	if paid+0.009 >= quote.Amount {
		if err := tx.SetStatus(requestID, status, string(models.CateringStatusPaid), actorType, actorID, note+"; paid in full"); err != nil {
			return nil, err
		}
		if err := tx.MarkPaidInFull(requestID, now); err != nil {
			return nil, err
		}
		receipt := &models.CateringReceipt{RequestID: requestID, ReceiptNumber: "RCT-" + request.Reference, TotalPaid: paid}
		if err := tx.InsertReceipt(receipt); err != nil {
			return nil, err
		}
		return &paymentOutcome{receipt: receipt}, nil
	}

	invoice := &models.CateringInvoice{
		RequestID: requestID, QuoteID: quote.ID, InvoiceNumber: *request.InvoiceNumber, Reason: "payment",
		Total: quote.Amount, DepositDue: depositFor(quote), AmountPaid: paid, Balance: roundNaira(quote.Amount - paid),
	}
	if err := tx.InsertInvoice(invoice); err != nil {
		return nil, err
	}
	return &paymentOutcome{invoice: invoice}, nil
}

func (s *cateringContractService) afterPayment(requestID int, outcome *paymentOutcome) (*models.CateringContract, error) {
	contract, err := s.GetContract(requestID)
	if err != nil {
		return nil, err
	}
	switch {
	case outcome.receipt != nil:
		s.notify(func() { s.sendReceiptEmail(contract, outcome.receipt) })
	case outcome.invoice != nil:
		s.notify(func() { s.sendInvoiceEmail(contract, outcome.invoice) })
	}
	return contract, nil
}

// StartOnlinePayment gives the customer a Paystack account to transfer the
// rest of the deposit or the whole balance to.
func (s *cateringContractService) StartOnlinePayment(input models.CateringPayOnlineInput) (*models.CateringPaymentCharge, error) {
	request, _, _, err := s.resolveCustomer(input.CateringCustomerAuth)
	if err != nil {
		return nil, err
	}
	contract, err := s.GetContract(request.ID)
	if err != nil {
		return nil, err
	}
	if contract.Account == nil || (request.Status != models.CateringStatusAccepted && request.Status != models.CateringStatusConfirmed) {
		return nil, invalid("There is no payment due on this booking")
	}
	if request.ContactEmail == nil || *request.ContactEmail == "" {
		return nil, invalid("Online payment needs an email address; please contact us to pay")
	}

	amount := contract.Account.Balance
	if input.Portion == "deposit" {
		amount = contract.Account.DepositRemaining
		if amount <= 0 {
			return nil, invalid("Your deposit is already paid")
		}
	}
	if amount <= 0 {
		return nil, invalid("Your booking is already paid in full")
	}

	if existing, err := s.contract.PendingCharge(request.ID, s.now()); err != nil {
		return nil, err
	} else if existing != nil && math.Abs(existing.Amount-amount) < 0.009 {
		return existing, nil
	}

	reference := CateringPaystackPrefix + strings.TrimPrefix(request.Reference, "CAT-") + "-" + randomSuffix()
	result, err := s.gateway.InitiateTransfer(*request.ContactEmail, amount, reference)
	if err != nil {
		return nil, fmt.Errorf("failed to start online payment: %w", err)
	}

	charge := &models.CateringPaymentCharge{
		RequestID:     request.ID,
		Reference:     result.Reference,
		Amount:        amount,
		AccountNumber: optionalString(result.AccountNumber),
		AccountName:   optionalString(result.AccountName),
		BankName:      optionalString(result.BankName),
		ExpiresAt:     result.ExpiresAt,
	}
	if err := s.contract.InTx(func(tx repository.CateringContractWriter) error {
		return tx.InsertCharge(charge)
	}); err != nil {
		return nil, err
	}
	return charge, nil
}

// adminTransitions are the status changes the admin makes by hand; the
// others (quoted, accepted, confirmed, paid) follow from quotes and payments.
var adminTransitions = map[string][]string{
	string(models.CateringStatusDeclined):  {"pending", "quoted"},
	string(models.CateringStatusCancelled): {"pending", "quoted", "accepted", "confirmed", "paid"},
	string(models.CateringStatusCompleted): {"paid"},
	string(models.CateringStatusPending):   {"declined", "cancelled"},
}

func (s *cateringContractService) ChangeStatus(requestID int, to string, adminID *int, note string) (*models.CateringContract, error) {
	allowedFrom, ok := adminTransitions[to]
	if !ok {
		return nil, invalid("%q is set automatically from quotes and payments", to)
	}
	request, err := s.request(requestID)
	if err != nil {
		return nil, err
	}

	err = s.contract.InTx(func(tx repository.CateringContractWriter) error {
		status, ok, err := tx.LockRequest(requestID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCateringRequestNotFound
		}
		permitted := false
		for _, from := range allowedFrom {
			if status == from {
				permitted = true
			}
		}
		if !permitted {
			return invalid("A %s request cannot be changed to %s", status, to)
		}
		if to == string(models.CateringStatusPending) && request.AcceptedQuoteID != nil {
			return invalid("This request had an accepted quote, so it cannot be reopened; create a new request instead")
		}
		if to == string(models.CateringStatusDeclined) || to == string(models.CateringStatusCancelled) {
			if err := tx.CloseOpenQuotes(requestID, models.QuoteStatusWithdrawn); err != nil {
				return err
			}
		}
		return tx.SetStatus(requestID, status, to, ActorAdmin, adminID, strings.TrimSpace(note))
	})
	if err != nil {
		return nil, err
	}
	return s.GetContract(requestID)
}

// GetContract loads a request with its quotes, ledger, invoices, receipt and
// history, as the admin sees it.
func (s *cateringContractService) GetContract(requestID int) (*models.CateringContract, error) {
	request, err := s.request(requestID)
	if err != nil {
		return nil, err
	}

	contract := &models.CateringContract{Request: request}
	if contract.Quotes, err = s.contract.ListQuotes(requestID); err != nil {
		return nil, err
	}
	for i := range contract.Quotes {
		contract.Quotes[i].DepositAmount = depositFor(&contract.Quotes[i])
		contract.Quotes[i].IsExpired = s.quoteExpired(&contract.Quotes[i])
	}
	contract.Quote = pickCurrentQuote(contract.Quotes, request.AcceptedQuoteID)

	if contract.Payments, err = s.contract.ListPayments(requestID); err != nil {
		return nil, err
	}
	if contract.Invoices, err = s.contract.ListInvoices(requestID); err != nil {
		return nil, err
	}
	if len(contract.Invoices) > 0 {
		contract.Invoice = &contract.Invoices[0]
	}
	if contract.Receipt, err = s.contract.GetReceipt(requestID); err != nil {
		return nil, err
	}
	if contract.History, err = s.contract.ListHistory(requestID); err != nil {
		return nil, err
	}
	if request.Status == models.CateringStatusAccepted || request.Status == models.CateringStatusConfirmed {
		if contract.PendingCharge, err = s.contract.PendingCharge(requestID, s.now()); err != nil {
			return nil, err
		}
	}

	if request.AcceptedQuoteID != nil && contract.Quote != nil {
		paid := 0.0
		for _, p := range contract.Payments {
			paid += p.Amount
		}
		paid = roundNaira(paid)
		deposit := depositFor(contract.Quote)
		contract.Account = &models.CateringAccount{
			Total:            contract.Quote.Amount,
			DepositDue:       deposit,
			AmountPaid:       paid,
			Balance:          math.Max(0, roundNaira(contract.Quote.Amount-paid)),
			DepositRemaining: math.Max(0, roundNaira(deposit-paid)),
		}
	}
	return contract, nil
}

// pickCurrentQuote is the accepted quote, else the newest open one.
func pickCurrentQuote(quotes []models.CateringQuote, acceptedID *int) *models.CateringQuote {
	for i := range quotes {
		if acceptedID != nil && quotes[i].ID == *acceptedID {
			return &quotes[i]
		}
	}
	if acceptedID == nil {
		for i := range quotes {
			if quotes[i].Status == models.QuoteStatusSent {
				return &quotes[i]
			}
		}
	}
	return nil
}

// CustomerView is the customer's view of their request and contract.
func (s *cateringContractService) CustomerView(auth models.CateringCustomerAuth) (*models.CateringContract, error) {
	request, _, _, err := s.resolveCustomer(auth)
	if err != nil {
		return nil, err
	}
	contract, err := s.GetContract(request.ID)
	if err != nil {
		return nil, err
	}
	return customerContract(contract), nil
}

// customerContract hides what only the business should see.
func customerContract(contract *models.CateringContract) *models.CateringContract {
	view := *contract
	request := *contract.Request
	request.AdminNotes = nil
	request.UserID = nil
	view.Request = &request
	view.Quotes = nil
	view.Invoices = nil
	view.History = nil
	view.Payments = make([]models.CateringPayment, 0, len(contract.Payments))
	for _, p := range contract.Payments {
		view.Payments = append(view.Payments, models.CateringPayment{
			ID: p.ID, RequestID: p.RequestID, Amount: p.Amount, Method: p.Method, PaidAt: p.PaidAt, CreatedAt: p.CreatedAt,
		})
	}
	return &view
}

// resolveCustomer finds the request for an emailed link token, or for a
// reference with the phone number it was made with.
func (s *cateringContractService) resolveCustomer(auth models.CateringCustomerAuth) (*models.CateringRequest, *models.CateringQuote, string, error) {
	if token := strings.TrimSpace(auth.Token); token != "" {
		quote, err := s.contract.QuoteByTokenHash(hashToken(token))
		if err != nil {
			return nil, nil, "", err
		}
		if quote == nil {
			return nil, nil, "", ErrCateringLinkInvalid
		}
		request, err := s.request(quote.RequestID)
		if err != nil {
			return nil, nil, "", err
		}
		return request, quote, models.AcceptedViaEmailLink, nil
	}

	if strings.TrimSpace(auth.Reference) == "" || strings.TrimSpace(auth.Phone) == "" {
		return nil, nil, "", invalid("Enter your reference and phone number")
	}
	request, err := s.requests.GetByReference(auth.Reference)
	if err != nil {
		return nil, nil, "", err
	}
	if request == nil || !utils.SamePhone(request.ContactPhone, auth.Phone) {
		return nil, nil, "", ErrCateringRequestNotFound
	}
	return request, nil, models.AcceptedViaStatusPage, nil
}

func (s *cateringContractService) currentQuote(requestID int) (*models.CateringQuote, error) {
	quotes, err := s.contract.ListQuotes(requestID)
	if err != nil {
		return nil, err
	}
	return pickCurrentQuote(quotes, nil), nil
}

func (s *cateringContractService) request(id int) (*models.CateringRequest, error) {
	request, err := s.requests.GetByID(id)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrCateringRequestNotFound
	}
	return request, nil
}

func (s *cateringContractService) notify(fn func()) {
	if s.email != nil {
		s.send(fn)
	}
}

// newQuoteToken returns a random link token and the hash stored for it, so
// the database never holds a usable link.
func newQuoteToken() (token, hash string, err error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("generating quote link: %w", err)
	}
	token = hex.EncodeToString(buf)
	return token, hashToken(token), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func randomSuffix() string {
	buf := make([]byte, 4)
	_, _ = rand.Read(buf)
	return strings.ToUpper(hex.EncodeToString(buf))
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
