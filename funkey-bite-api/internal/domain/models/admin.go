package models

import "time"

type AdminStats struct {
	TotalOrders    int             `json:"totalOrders"`
	TotalRevenue   float64         `json:"totalRevenue"`
	PendingOrders  int             `json:"pendingOrders"`
	ActiveCatering int             `json:"activeCatering"`
	NewCustomers   int             `json:"newCustomers"`
	PopularItems   []MenuItemStats `json:"popularItems"`
}

type MenuItemStats struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	TotalSold int     `json:"totalSold"`
	Revenue   float64 `json:"revenue"`
}

type SalesReport struct {
	Date         string  `json:"date"`
	TotalOrders  int     `json:"totalOrders"`
	TotalRevenue float64 `json:"totalRevenue"`
	AverageOrder float64 `json:"averageOrder"`
}

type AdminUser struct {
	ID           int        `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         string     `json:"role" db:"role"`
	IsActive     bool       `json:"isActive" db:"is_active"`
	LastLogin    *time.Time `json:"lastLogin,omitempty" db:"last_login"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
}

type BusinessSettings struct {
	ID             int            `json:"id" db:"id"`
	BusinessName   string         `json:"businessName" db:"business_name"`
	PhoneNumber    string         `json:"phoneNumber" db:"phone_number"`
	Email          string         `json:"email" db:"email"`
	Address        string         `json:"address" db:"address"` // street line
	City           string         `json:"city" db:"city"`
	State          string         `json:"state" db:"state"`
	PostalCode     string         `json:"postalCode" db:"postal_code"`
	Country        string         `json:"country" db:"country"`
	OpeningHours   []OpeningHours `json:"openingHours" db:"opening_hours"`
	DeliveryFee    float64        `json:"deliveryFee" db:"delivery_fee"`
	MinOrderAmount float64        `json:"minOrderAmount" db:"min_order_amount"`
	TaxRate        float64        `json:"taxRate" db:"tax_rate"`
	IsDeliveryOpen bool           `json:"isDeliveryOpen" db:"is_delivery_open"`
	IsPickupOpen   bool           `json:"isPickupOpen" db:"is_pickup_open"`
	// Timezone is the IANA zone the opening hours are in (BUSINESS_TIMEZONE);
	// it is configuration, not a column.
	Timezone  string    `json:"timezone" db:"-"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// OpeningHours is one weekday's trading window. Times are 24-hour "HH:MM" in
// the business timezone; a CloseTime at or before OpenTime means the business
// closes after midnight (e.g. 18:00-02:00, or 08:00-00:00 for midnight).
type OpeningHours struct {
	Day       string `json:"day"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	IsOpen    bool   `json:"isOpen"`
}
