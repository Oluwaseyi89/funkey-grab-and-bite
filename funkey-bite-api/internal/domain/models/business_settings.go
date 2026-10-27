package models

type BusinessSettingsUpdate struct {
	BusinessName         *string         `json:"businessName,omitempty"`
	PhoneNumber          *string         `json:"phoneNumber,omitempty"`
	Email                *string         `json:"email,omitempty"`
	Address              *string         `json:"address,omitempty"`
	City                 *string         `json:"city,omitempty"`
	State                *string         `json:"state,omitempty"`
	PostalCode           *string         `json:"postalCode,omitempty"`
	Country              *string         `json:"country,omitempty"`
	OpeningHours         *[]OpeningHours `json:"openingHours,omitempty"`
	DeliveryFee          *float64        `json:"deliveryFee,omitempty"`
	MinOrderAmount       *float64        `json:"minOrderAmount,omitempty"`
	TaxRate              *float64        `json:"taxRate,omitempty"`
	IsDeliveryOpen       *bool           `json:"isDeliveryOpen,omitempty"`
	IsPickupOpen         *bool           `json:"isPickupOpen,omitempty"`
	PayOnDeliveryEnabled *bool           `json:"payOnDeliveryEnabled,omitempty"`
	CateringTerms        *string         `json:"cateringTerms,omitempty"`
}
