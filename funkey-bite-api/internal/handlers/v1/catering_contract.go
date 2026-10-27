package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/realtime"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

// CateringContractHandler serves quotes, acceptance and payments for
// catering requests.
type CateringContractHandler struct {
	contracts services.CateringContractService
	validate  *validator.Validate
}

func NewCateringContractHandler(contracts services.CateringContractService, catering *CateringHandler) *CateringContractHandler {
	return &CateringContractHandler{contracts: contracts, validate: catering.validate}
}

var contractFieldLabels = map[string]string{
	"amount": "Amount", "lineItems": "Quote lines", "description": "Line description",
	"depositPercent": "Deposit percentage", "terms": "Terms", "message": "Message",
	"validDays": "Days the quote is valid", "method": "Payment method", "reference": "Reference",
	"note": "Note", "portion": "Amount to pay", "token": "Link", "phone": "Phone number",
}

func respondContractError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, services.ErrCateringLinkInvalid) {
		c.JSON(http.StatusNotFound, gin.H{"error": "This link is no longer valid. Look up your request with its reference and phone number instead."})
		return
	}
	respondCateringError(c, err, fallback)
}

func adminActor(c *gin.Context) *int {
	if id, ok := c.Get("user_id"); ok {
		if adminID, ok := id.(int); ok {
			return &adminID
		}
	}
	return nil
}

func (h *CateringContractHandler) bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return false
	}
	if err := h.validate.Struct(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationMessage(err, contractFieldLabels)})
		return false
	}
	return true
}

func requestID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return 0, false
	}
	return id, true
}

func (h *CateringContractHandler) broadcast(contract *models.CateringContract) {
	realtime.GlobalBroadcaster.Broadcast("catering_request_updated", contract.Request)
}

// HandlePaystackCharge records a catering payment from the Paystack webhook
// and tells the admin app.
func (h *CateringContractHandler) HandlePaystackCharge(reference string, amountKobo int64) (bool, error) {
	handled, err := h.contracts.HandlePaystackCharge(reference, amountKobo)
	if handled && err == nil {
		realtime.GlobalBroadcaster.Broadcast("catering_request_updated", gin.H{"paymentReference": reference})
	}
	return handled, err
}

// GetContract returns a request with its quotes, payments, invoices,
// receipt and history
// @Summary Get catering contract
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Success 200 {object} models.CateringContract
// @Router /admin/catering/requests/{id}/contract [get]
func (h *CateringContractHandler) GetContract(c *gin.Context) {
	id, ok := requestID(c)
	if !ok {
		return
	}
	contract, err := h.contracts.GetContract(id)
	if err != nil {
		respondContractError(c, err, "Failed to fetch catering request")
		return
	}
	c.JSON(http.StatusOK, contract)
}

// SendQuote sends the customer a quote, replacing any open one
// @Summary Send catering quote
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Param quote body models.CateringQuoteInput true "Quote"
// @Success 201 {object} models.CateringContract
// @Router /admin/catering/requests/{id}/quotes [post]
func (h *CateringContractHandler) SendQuote(c *gin.Context) {
	id, ok := requestID(c)
	if !ok {
		return
	}
	var input models.CateringQuoteInput
	if !h.bind(c, &input) {
		return
	}
	contract, err := h.contracts.SendQuote(id, &input, adminActor(c))
	if err != nil {
		respondContractError(c, err, "Failed to send quote")
		return
	}
	h.broadcast(contract)
	c.JSON(http.StatusCreated, contract)
}

// RecordPayment records a cash or bank transfer payment
// @Summary Record catering payment
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Param payment body models.CateringPaymentInput true "Payment"
// @Success 201 {object} models.CateringContract
// @Router /admin/catering/requests/{id}/payments [post]
func (h *CateringContractHandler) RecordPayment(c *gin.Context) {
	id, ok := requestID(c)
	if !ok {
		return
	}
	var input models.CateringPaymentInput
	if !h.bind(c, &input) {
		return
	}
	contract, err := h.contracts.RecordPayment(id, &input, adminActor(c))
	if err != nil {
		respondContractError(c, err, "Failed to record payment")
		return
	}
	h.broadcast(contract)
	c.JSON(http.StatusCreated, contract)
}

type cateringStatusChange struct {
	Status string `json:"status" validate:"required"`
	Note   string `json:"note" validate:"max=500"`
}

// ChangeStatus declines, cancels, completes or reopens a request
// @Summary Change catering request status
// @Description Only declined, cancelled, completed and pending can be set by hand; the rest follow from quotes and payments
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Catering Request ID"
// @Success 200 {object} models.CateringContract
// @Router /admin/catering/requests/{id}/status [patch]
func (h *CateringContractHandler) ChangeStatus(c *gin.Context) {
	id, ok := requestID(c)
	if !ok {
		return
	}
	var body cateringStatusChange
	if !h.bind(c, &body) {
		return
	}
	contract, err := h.contracts.ChangeStatus(id, body.Status, adminActor(c), body.Note)
	if err != nil {
		respondContractError(c, err, "Failed to update status")
		return
	}
	h.broadcast(contract)
	c.JSON(http.StatusOK, contract)
}

// ViewBooking shows the customer their request, quote and payments
// @Summary View catering booking
// @Description Identify with the emailed link token, or the reference and phone number
// @Tags catering
// @Accept json
// @Produce json
// @Param auth body models.CateringCustomerAuth true "Link token, or reference and phone"
// @Success 200 {object} models.CateringContract
// @Router /catering/booking [post]
func (h *CateringContractHandler) ViewBooking(c *gin.Context) {
	var auth models.CateringCustomerAuth
	if !h.bind(c, &auth) {
		return
	}
	contract, err := h.contracts.CustomerView(auth)
	if err != nil {
		respondContractError(c, err, "Failed to fetch your booking")
		return
	}
	c.JSON(http.StatusOK, contract)
}

// AcceptQuote accepts the open quote and issues the invoice
// @Summary Accept catering quote
// @Tags catering
// @Accept json
// @Produce json
// @Param accept body models.CateringAcceptInput true "Identification and agreement to the terms"
// @Success 200 {object} models.CateringContract
// @Router /catering/booking/accept [post]
func (h *CateringContractHandler) AcceptQuote(c *gin.Context) {
	var input models.CateringAcceptInput
	if !h.bind(c, &input) {
		return
	}
	contract, err := h.contracts.AcceptQuote(input)
	if err != nil {
		respondContractError(c, err, "Failed to accept the quote")
		return
	}
	h.broadcast(contract)
	c.JSON(http.StatusOK, contract)
}

// PayOnline gives the customer an account to transfer the deposit or balance to
// @Summary Pay for catering online
// @Tags catering
// @Accept json
// @Produce json
// @Param pay body models.CateringPayOnlineInput true "Identification and what to pay"
// @Success 200 {object} models.CateringPaymentCharge
// @Router /catering/booking/pay [post]
func (h *CateringContractHandler) PayOnline(c *gin.Context) {
	var input models.CateringPayOnlineInput
	if !h.bind(c, &input) {
		return
	}
	charge, err := h.contracts.StartOnlinePayment(input)
	if err != nil {
		respondContractError(c, err, "Failed to start the payment")
		return
	}
	c.JSON(http.StatusOK, charge)
}
