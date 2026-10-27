package v1

import (
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

func TestCateringValidationMessagesAreReadable(t *testing.T) {
	handler := NewCateringHandler(nil)

	cases := []struct {
		name  string
		input models.CateringRequestInput
		want  string
	}{
		{"missing type", models.CateringRequestInput{}, "Request type is required"},
		{"bad email", func() models.CateringRequestInput {
			in := validCateringInput()
			in.ContactEmail = "not-an-email"
			return in
		}(), "Enter a valid email address"},
		{"unknown event type", func() models.CateringRequestInput {
			in := validCateringInput()
			in.EventType = "rave"
			return in
		}(), "Event type has an option we don't offer"},
		{"unknown dietary option", func() models.CateringRequestInput {
			in := validCateringInput()
			in.DietaryRequirements = []string{"keto"}
			return in
		}(), "Dietary requirements has an option we don't offer"},
		{"no venue", func() models.CateringRequestInput {
			in := validCateringInput()
			in.VenueAddress = ""
			return in
		}(), "Venue address is required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := handler.validate.Struct(tc.input)
			if got := validationMessage(err); got != tc.want {
				t.Fatalf("validationMessage() = %q, want %q", got, tc.want)
			}
		})
	}

	if err := handler.validate.Struct(validCateringInput()); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}
}

func validCateringInput() models.CateringRequestInput {
	return models.CateringRequestInput{
		RequestType: "custom", ContactName: "Ada Obi", ContactPhone: "08012345678",
		ContactEmail: "ada@example.com", EventType: "party", EventDate: "2026-11-10",
		EventTime: "14:00", DurationHours: 4, GuestCount: 60, VenueAddress: "5 Bourdillon Road",
		VenueArea: "Ikoyi", VenueState: "Lagos", ServiceStyle: "buffet",
		DietaryRequirements: []string{"halal"},
	}
}
