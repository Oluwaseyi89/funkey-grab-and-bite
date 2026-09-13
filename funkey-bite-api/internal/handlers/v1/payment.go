package v1

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// Webhook receives Paystack's server-to-server payment notifications. It must
// read the raw request body (not ShouldBindJSON) because the signature is
// computed over the exact bytes Paystack sent.
func (h *PaymentHandler) Webhook(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		handlers.Error(c, http.StatusBadRequest, "INVALID_BODY", "Failed to read webhook body")
		return
	}

	signature := c.GetHeader("x-paystack-signature")

	if err := h.paymentService.HandleWebhook(rawBody, signature); err != nil {
		if errors.Is(err, services.ErrInvalidWebhookSignature) {
			handlers.Error(c, http.StatusUnauthorized, "INVALID_SIGNATURE", "Webhook signature verification failed")
			return
		}
		// Any other failure gets a 500 so Paystack retries delivery — the failure
		// modes here (DB errors, a not-yet-visible order on a replica) are the
		// kind that self-resolve on retry.
		handlers.Error(c, http.StatusInternalServerError, "WEBHOOK_PROCESSING_FAILED", err.Error())
		return
	}

	handlers.Success(c, gin.H{"received": true})
}
