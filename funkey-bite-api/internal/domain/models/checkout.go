package models

// GuestCustomer is someone checking out without an account and without a
// customer record yet. The record is created with the order, matched by phone.
type GuestCustomer struct {
	Name  string
	Phone string
	Email *string
}

// Checkout is everything one order saves. It is saved in a single
// transaction: the order, its items, the stock taken, the customer record and
// the promotion use are all saved together or not at all.
type Checkout struct {
	Order *Order
	Items []OrderItemRequest
	// Guest is set when Order.UserID is nil; their record is found or created
	// by phone and linked to the order.
	Guest *GuestCustomer
	// Promotion, when set, records the use of PromotionID with
	// DiscountApplied. The order and customer are filled in when saving.
	Promotion *PromotionUsage
}
