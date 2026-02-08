package models

// type BusinessSettings struct {
// 	ID             int       `json:"id" db:"id"`
// 	BusinessName   string    `json:"businessName" db:"business_name"`
// 	PhoneNumber    string    `json:"phoneNumber" db:"phone_number"`
// 	Email          string    `json:"email" db:"email"`
// 	Address        string    `json:"address" db:"address"`
// 	OpeningHours   string    `json:"openingHours" db:"opening_hours"` // JSON string
// 	DeliveryFee    float64   `json:"deliveryFee" db:"delivery_fee"`
// 	MinOrderAmount float64   `json:"minOrderAmount" db:"min_order_amount"`
// 	TaxRate        float64   `json:"taxRate" db:"tax_rate"`
// 	IsDeliveryOpen bool      `json:"isDeliveryOpen" db:"is_delivery_open"`
// 	IsPickupOpen   bool      `json:"isPickupOpen" db:"is_pickup_open"`
// 	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
// 	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
// }

type BusinessHours struct {
	Day       string `json:"day"`       // Monday, Tuesday, etc.
	OpenTime  string `json:"openTime"`  // 08:00
	CloseTime string `json:"closeTime"` // 22:00
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
