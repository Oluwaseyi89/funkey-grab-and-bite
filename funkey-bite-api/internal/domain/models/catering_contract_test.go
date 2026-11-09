package models

import (
	"encoding/json"
	"testing"
)

// The booking reference and the payment reference travel in one body; they
// must not overwrite each other.
func TestVerifyPaymentInputKeepsBothReferences(t *testing.T) {
	var input CateringVerifyPaymentInput
	if err := json.Unmarshal([]byte(`{"reference":"CAT-ABC123","phone":"08012345678","paymentReference":"CATPAY-ABC123-1"}`), &input); err != nil {
		t.Fatal(err)
	}
	if input.Reference != "CAT-ABC123" || input.PaymentReference != "CATPAY-ABC123-1" {
		t.Fatalf("expected both references kept, got %+v", input)
	}
}
