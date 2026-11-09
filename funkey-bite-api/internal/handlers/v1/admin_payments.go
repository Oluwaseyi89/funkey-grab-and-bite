package v1

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers/middleware"
	"funkey-grab-and-bite/funkey-bite-api/internal/realtime"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type onlinePaymentReader interface {
	List(ctx context.Context, filter models.OnlinePaymentFilter, limit, offset int) ([]models.OnlinePayment, int, error)
	CountByStatus(ctx context.Context) (map[string]int, error)
	Unmatched(ctx context.Context, limit int) ([]models.PaymentEventRecord, error)
}

type paymentReviewer interface {
	ResolvePaymentReview(ctx context.Context, orderID int, decision, note string) (string, error)
}

type jobAdmin interface {
	List(ctx context.Context, status string, limit, offset int) ([]models.JobSummary, int, error)
	RetryDead(ctx context.Context, id int64) (bool, error)
}

type auditReader interface {
	List(ctx context.Context, filter models.AuditFilter, limit, offset int) ([]models.AuditRecord, int, error)
}

// AdminPaymentsHandler serves the admin's payments, review queue, failed
// jobs and audit log.
type AdminPaymentsHandler struct {
	payments onlinePaymentReader
	reviewer paymentReviewer
	checker  services.PaymentService
	jobs     jobAdmin
	audit    auditReader
}

func NewAdminPaymentsHandler(payments onlinePaymentReader, reviewer paymentReviewer, checker services.PaymentService, jobs jobAdmin, audit auditReader) *AdminPaymentsHandler {
	return &AdminPaymentsHandler{payments: payments, reviewer: reviewer, checker: checker, jobs: jobs, audit: audit}
}

func pageParams(c *gin.Context) (page, limit, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit, (page - 1) * limit
}

func pagination(page, limit, total int) gin.H {
	return gin.H{"page": page, "limit": limit, "total": total, "totalPages": (total + limit - 1) / limit}
}

// ListPayments returns online transfer orders with counts per payment status
// @Summary List online payments
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param status query string false "Payment status, e.g. needs_review"
// @Param search query string false "Order number, customer, phone or reference"
// @Router /admin/payments [get]
func (h *AdminPaymentsHandler) ListPayments(c *gin.Context) {
	page, limit, offset := pageParams(c)
	filter := models.OnlinePaymentFilter{
		PaymentStatus: c.Query("status"),
		Search:        strings.TrimSpace(c.Query("search")),
	}
	payments, total, err := h.payments.List(c.Request.Context(), filter, limit, offset)
	if err != nil {
		log.Printf("List payments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load payments"})
		return
	}
	counts, err := h.payments.CountByStatus(c.Request.Context())
	if err != nil {
		log.Printf("Count payments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load payments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payments": payments, "counts": counts, "pagination": pagination(page, limit, total)})
}

// ListUnmatchedPayments returns payments that matched no order
// @Summary List unmatched payments
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Router /admin/payments/unmatched [get]
func (h *AdminPaymentsHandler) ListUnmatchedPayments(c *gin.Context) {
	events, err := h.payments.Unmatched(c.Request.Context(), 50)
	if err != nil {
		log.Printf("List unmatched payments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load unmatched payments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payments": events})
}

type paymentReviewRequest struct {
	Decision string `json:"decision"`
	Note     string `json:"note"`
}

// ResolvePaymentReview accepts a payment, marks it for refund, or records the
// refund as made
// @Summary Decide on a payment
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Router /admin/orders/{id}/payment-review [post]
func (h *AdminPaymentsHandler) ResolvePaymentReview(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	var req paymentReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	switch req.Decision {
	case models.ReviewDecisionAccept, models.ReviewDecisionRefund, models.ReviewDecisionRefunded:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Decision must be accept, refund or refunded"})
		return
	}
	note := strings.TrimSpace(req.Note)
	if len(note) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note must be 500 characters or fewer"})
		return
	}

	paymentStatus, err := h.reviewer.ResolvePaymentReview(middleware.RequestContext(c), orderID, req.Decision, note)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	case errors.Is(err, repository.ErrNotUnderReview), errors.Is(err, repository.ErrCannotReopen):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		log.Printf("Resolve payment review for order %d: %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record the decision"})
		return
	}

	realtime.GlobalBroadcaster.Broadcast("order_updated", gin.H{"id": orderID, "paymentStatus": paymentStatus})
	c.JSON(http.StatusOK, gin.H{"orderId": orderID, "paymentStatus": paymentStatus})
}

// RecheckPayment asks Paystack whether an order was paid
// @Summary Re-check a payment with Paystack
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Router /admin/orders/{id}/payment-check [post]
func (h *AdminPaymentsHandler) RecheckPayment(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	result, paid, err := h.checker.RecheckOnlinePayment(middleware.RequestContext(c), orderID)
	switch {
	case errors.Is(err, services.ErrOrderNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	case errors.Is(err, services.ErrNotOnlinePayment):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		log.Printf("Re-check payment for order %d: %v", orderID, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Couldn't reach Paystack. Please try again."})
		return
	}

	message := "Paystack has no payment for this order yet."
	if paid {
		switch result.Outcome {
		case repository.PaymentApplied:
			message = "Paystack confirmed the payment; the order is now paid."
			realtime.GlobalBroadcaster.Broadcast("order_updated", gin.H{"id": orderID, "paymentStatus": models.PaymentStatusPaid})
		case repository.PaymentNeedsReview:
			message = "Paystack has a payment that needs your review."
		default:
			message = "The payment was already recorded."
		}
	}
	c.JSON(http.StatusOK, gin.H{"orderId": orderID, "paid": paid, "outcome": result.Outcome, "message": message})
}

// ListJobs returns background jobs, failed ones by default
// @Summary List background jobs
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param status query string false "pending, running, succeeded, dead (default) or all"
// @Router /admin/jobs [get]
func (h *AdminPaymentsHandler) ListJobs(c *gin.Context) {
	page, limit, offset := pageParams(c)
	status := c.DefaultQuery("status", "dead")
	switch status {
	case "all":
		status = ""
	case "pending", "running", "succeeded", "dead":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status filter"})
		return
	}
	jobs, total, err := h.jobs.List(c.Request.Context(), status, limit, offset)
	if err != nil {
		log.Printf("List jobs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs, "pagination": pagination(page, limit, total)})
}

// RetryJob gives a failed job a fresh set of attempts
// @Summary Retry a failed job
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Router /admin/jobs/{id}/retry [post]
func (h *AdminPaymentsHandler) RetryJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	retried, err := h.jobs.RetryDead(c.Request.Context(), id)
	if err != nil {
		log.Printf("Retry job %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry the job"})
		return
	}
	if !retried {
		c.JSON(http.StatusConflict, gin.H{"error": "Only failed jobs can be retried"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "pending"})
}

// ListAudit returns audit log entries, newest first
// @Summary List audit log entries
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Router /admin/audit [get]
func (h *AdminPaymentsHandler) ListAudit(c *gin.Context) {
	page, limit, offset := pageParams(c)
	filter := models.AuditFilter{
		EntityType: c.Query("entityType"),
		EntityID:   c.Query("entityId"),
		Action:     c.Query("action"),
		ActorType:  c.Query("actorType"),
	}
	entries, total, err := h.audit.List(c.Request.Context(), filter, limit, offset)
	if err != nil {
		log.Printf("List audit entries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load the audit log"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"entries": entries, "pagination": pagination(page, limit, total)})
}
