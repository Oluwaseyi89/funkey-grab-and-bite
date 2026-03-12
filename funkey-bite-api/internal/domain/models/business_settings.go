package models

type BusinessHours struct {
	Day       string `json:"day"`       
	OpenTime  string `json:"openTime"` 
	CloseTime string `json:"closeTime"` 
	IsOpen    bool   `json:"isOpen"`
}

type BusinessSettingsUpdate struct {
	BusinessName   *string  `json:"businessName,omitempty"`
	PhoneNumber    *string  `json:"phoneNumber,omitempty"`
	Email          *string  `json:"email,omitempty"`
	Address        *string  `json:"address,omitempty"`
	OpeningHours   *string  `json:"openingHours,omitempty"`
	DeliveryFee    *float64 `json:"deliveryFee,omitempty"`
	MinOrderAmount *float64 `json:"minOrderAmount,omitempty"`
	TaxRate        *float64 `json:"taxRate,omitempty"`
	IsDeliveryOpen *bool    `json:"isDeliveryOpen,omitempty"`
	IsPickupOpen   *bool    `json:"isPickupOpen,omitempty"`
}
