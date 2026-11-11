package models

// Sales figures for the admin dashboard and reports. Revenue is money
// received: orders whose payment is recorded as paid and that weren't
// cancelled, counted on the day the payment was made. Orders are those
// placed in the period, cancelled ones left out. Days are calendar days in
// the business timezone.

// PeriodTotals sums up a stretch of time.
type PeriodTotals struct {
	Revenue      float64 `json:"revenue"`
	Orders       int     `json:"orders"`
	PaidOrders   int     `json:"paidOrders"`
	AverageOrder float64 `json:"averageOrder"`
	Cancelled    int     `json:"cancelled"`
}

// SalesPoint is one day, week or month. Periods with no sales are included
// with zeros, so charts don't skip them.
type SalesPoint struct {
	// Start of the period, YYYY-MM-DD (weeks start on Monday).
	Date       string  `json:"date"`
	Revenue    float64 `json:"revenue"`
	Orders     int     `json:"orders"`
	PaidOrders int     `json:"paidOrders"`
}

// HourPoint counts orders placed in one hour of the day (0-23).
type HourPoint struct {
	Hour   int `json:"hour"`
	Orders int `json:"orders"`
}

// TopItem is a best seller.
type TopItem struct {
	MenuItemID int     `json:"menuItemId"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Revenue    float64 `json:"revenue"`
}

// Share is part of a whole, e.g. revenue by payment method.
type Share struct {
	Key     string  `json:"key"`
	Orders  int     `json:"orders"`
	Revenue float64 `json:"revenue"`
}

// OrderQueue counts orders that still need the kitchen.
type OrderQueue struct {
	Pending   int `json:"pending"`
	Confirmed int `json:"confirmed"`
	Preparing int `json:"preparing"`
	Ready     int `json:"ready"`
}

// DashboardOverview is what the admin sees first: today so far, what needs
// doing, and the last two weeks.
type DashboardOverview struct {
	Timezone string `json:"timezone"`
	// Today up to now, and yesterday up to the same time, for a fair comparison.
	Today          PeriodTotals `json:"today"`
	YesterdaySoFar PeriodTotals `json:"yesterdaySoFar"`
	Queue          OrderQueue   `json:"queue"`
	// Online payments that arrived but need checking, and ones still awaited.
	PaymentsToReview int `json:"paymentsToReview"`
	AwaitingPayment  int `json:"awaitingPayment"`
	LowStock         int `json:"lowStock"`
	// New catering requests waiting for a quote.
	CateringToQuote int `json:"cateringToQuote"`
	// The last 14 days including today, oldest first.
	Daily []SalesPoint `json:"daily"`
	// Orders by hour today.
	Hourly []HourPoint `json:"hourly"`
	// Best sellers over the last 7 days.
	TopItems []TopItem `json:"topItems"`
}

// SalesSummary is the reports page for a chosen period.
type SalesSummary struct {
	From     string       `json:"from"`
	To       string       `json:"to"`
	Group    string       `json:"group"`
	Timezone string       `json:"timezone"`
	Totals   PeriodTotals `json:"totals"`
	// The same number of days immediately before From.
	Previous       PeriodTotals `json:"previous"`
	PreviousFrom   string       `json:"previousFrom"`
	PreviousTo     string       `json:"previousTo"`
	Series         []SalesPoint `json:"series"`
	Hourly         []HourPoint  `json:"hourly"`
	TopItems       []TopItem    `json:"topItems"`
	PaymentMethods []Share      `json:"paymentMethods"`
	OrderTypes     []Share      `json:"orderTypes"`
}
