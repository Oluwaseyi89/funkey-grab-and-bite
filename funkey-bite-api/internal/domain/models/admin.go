package models

import "time"

// AdminStats represents dashboard statistics
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

// SalesReport represents sales data for reporting
type SalesReport struct {
	Date         string  `json:"date"`
	TotalOrders  int     `json:"totalOrders"`
	TotalRevenue float64 `json:"totalRevenue"`
	AverageOrder float64 `json:"averageOrder"`
}

// AdminUser represents admin users for the dashboard
type AdminUser struct {
	ID           int        `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         string     `json:"role" db:"role"` // admin, manager, staff
	IsActive     bool       `json:"isActive" db:"is_active"`
	LastLogin    *time.Time `json:"lastLogin,omitempty" db:"last_login"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
}

// BusinessSettings represents business configuration
type BusinessSettings struct {
	BusinessName   string         `json:"businessName" db:"business_name"`
	PhoneNumber    string         `json:"phoneNumber" db:"phone_number"`
	Email          string         `json:"email" db:"email"`
	Address        string         `json:"address" db:"address"`
	OpeningHours   []OpeningHours `json:"openingHours" db:"opening_hours"`
	DeliveryFee    float64        `json:"deliveryFee" db:"delivery_fee"`
	MinOrderAmount float64        `json:"minOrderAmount" db:"min_order_amount"`
	TaxRate        float64        `json:"taxRate" db:"tax_rate"`
}

type OpeningHours struct {
	Day       string `json:"day"`       // Monday, Tuesday, etc.
	OpenTime  string `json:"openTime"`  // 08:00
	CloseTime string `json:"closeTime"` // 22:00
	IsOpen    bool   `json:"isOpen"`
}
