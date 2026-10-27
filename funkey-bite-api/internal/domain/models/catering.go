package models

import "time"

type CateringStatus string

// A request moves pending → quoted → confirmed → completed; it can be
// declined by the business or cancelled by the customer before it happens.
const (
	CateringStatusPending   CateringStatus = "pending"
	CateringStatusQuoted    CateringStatus = "quoted"
	CateringStatusConfirmed CateringStatus = "confirmed"
	CateringStatusDeclined  CateringStatus = "declined"
	CateringStatusCompleted CateringStatus = "completed"
	CateringStatusCancelled CateringStatus = "cancelled"
)

// IsValidCateringStatus reports whether status is a known request status.
func IsValidCateringStatus(status string) bool {
	switch CateringStatus(status) {
	case CateringStatusPending, CateringStatusQuoted, CateringStatusConfirmed,
		CateringStatusDeclined, CateringStatusCompleted, CateringStatusCancelled:
		return true
	}
	return false
}

// Catering request types: for one of the admin's packages, or custom.
const (
	CateringRequestTypePackage = "package"
	CateringRequestTypeCustom  = "custom"
)

// How food is served at the event.
const (
	ServiceStyleDropOff     = "drop_off"     // food delivered ready to serve
	ServiceStyleBuffet      = "buffet"       // buffet set up and attended by staff
	ServiceStyleFullService = "full_service" // served by staff, setup and cleanup
)

// CateringEventTypes are the occasions a request can be for; "other" needs a
// short description.
var CateringEventTypes = []string{
	"wedding", "party", "birthday", "naming", "burial", "church",
	"corporate", "graduation", "other",
}

// CateringDietaryOptions are the dietary needs a request can flag.
var CateringDietaryOptions = []string{
	"vegetarian", "vegan", "halal", "no_pork", "nut_free", "gluten_free", "dairy_free",
}

// CateringPackage is an admin-defined catering offer.
type CateringPackage struct {
	ID             int                   `json:"id"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	PricePerPerson float64               `json:"pricePerPerson"`
	MinGuests      int                   `json:"minGuests"`
	MaxGuests      *int                  `json:"maxGuests,omitempty"`
	ServiceStyle   string                `json:"serviceStyle"`
	MinNoticeDays  int                   `json:"minNoticeDays"`
	Inclusions     []string              `json:"inclusions"`
	DishChoices    int                   `json:"dishChoices"`
	MenuItems      []CateringPackageItem `json:"menuItems"`
	IsActive       bool                  `json:"isActive"`
	DisplayOrder   int                   `json:"displayOrder"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
}

// CateringPackageItem is a menu item customers may choose for a package.
type CateringPackageItem struct {
	MenuItemID  int    `json:"menuItemId"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"isAvailable"`
}

// CateringPackageInput creates or replaces a package.
type CateringPackageInput struct {
	Name           string   `json:"name" validate:"required,min=2,max=120"`
	Description    string   `json:"description" validate:"max=1000"`
	PricePerPerson float64  `json:"pricePerPerson" validate:"min=0"`
	MinGuests      int      `json:"minGuests" validate:"required,min=1,max=1000"`
	MaxGuests      *int     `json:"maxGuests,omitempty" validate:"omitempty,min=1,max=1000"`
	ServiceStyle   string   `json:"serviceStyle" validate:"required,oneof=drop_off buffet full_service"`
	MinNoticeDays  int      `json:"minNoticeDays" validate:"min=0,max=180"`
	Inclusions     []string `json:"inclusions" validate:"max=20,dive,required,max=120"`
	DishChoices    int      `json:"dishChoices" validate:"min=0,max=20"`
	MenuItemIDs    []int    `json:"menuItemIds" validate:"max=60,dive,min=1"`
	IsActive       bool     `json:"isActive"`
	DisplayOrder   int      `json:"displayOrder" validate:"min=0,max=1000"`
}

// CateringSelectedItem records a dish chosen in a request, with its name at
// the time so the request still reads correctly if the menu changes.
type CateringSelectedItem struct {
	MenuItemID int    `json:"menuItemId"`
	Name       string `json:"name"`
}

type CateringRequest struct {
	ID                  int                    `json:"id"`
	Reference           string                 `json:"reference"`
	RequestType         string                 `json:"requestType"`
	UserID              *int                   `json:"userId,omitempty"`
	PackageID           *int                   `json:"packageId,omitempty"`
	PackageName         *string                `json:"packageName,omitempty"`
	ContactName         string                 `json:"contactName"`
	ContactPhone        string                 `json:"contactPhone"`
	ContactEmail        *string                `json:"contactEmail,omitempty"`
	Organization        *string                `json:"organization,omitempty"`
	EventType           string                 `json:"eventType"`
	EventTypeOther      *string                `json:"eventTypeOther,omitempty"`
	EventName           *string                `json:"eventName,omitempty"`
	EventDate           string                 `json:"eventDate"`
	EventTime           *string                `json:"eventTime,omitempty"`
	DurationHours       *float64               `json:"durationHours,omitempty"`
	GuestCount          int                    `json:"guestCount"`
	VenueName           *string                `json:"venueName,omitempty"`
	VenueAddress        *string                `json:"venueAddress,omitempty"`
	VenueArea           *string                `json:"venueArea,omitempty"`
	VenueState          *string                `json:"venueState,omitempty"`
	VenueNotes          *string                `json:"venueNotes,omitempty"`
	ServiceStyle        *string                `json:"serviceStyle,omitempty"`
	SelectedItems       []CateringSelectedItem `json:"selectedItems"`
	MenuNotes           *string                `json:"menuNotes,omitempty"`
	DietaryRequirements []string               `json:"dietaryRequirements"`
	DietaryNotes        *string                `json:"dietaryNotes,omitempty"`
	Budget              *float64               `json:"budget,omitempty"`
	SpecialRequests     *string                `json:"specialRequests,omitempty"`
	EstimatedTotal      *float64               `json:"estimatedTotal,omitempty"`
	QuotedAmount        *float64               `json:"quotedAmount,omitempty"`
	AdminNotes          *string                `json:"adminNotes,omitempty"`
	Status              CateringStatus         `json:"status"`
	CreatedAt           time.Time              `json:"createdAt"`
	UpdatedAt           *time.Time             `json:"updatedAt,omitempty"`
}

// CateringRequestInput is a customer's catering request. Rules that depend
// on the request type or the package are checked by the catering service.
type CateringRequestInput struct {
	RequestType         string   `json:"requestType" validate:"required,oneof=package custom"`
	PackageID           *int     `json:"packageId,omitempty" validate:"omitempty,min=1"`
	ContactName         string   `json:"contactName" validate:"required,min=2,max=200"`
	ContactPhone        string   `json:"contactPhone" validate:"required,max=30"`
	ContactEmail        string   `json:"contactEmail" validate:"required,email,max=200"`
	Organization        *string  `json:"organization,omitempty" validate:"omitempty,max=200"`
	EventType           string   `json:"eventType" validate:"required,oneof=wedding party birthday naming burial church corporate graduation other"`
	EventTypeOther      *string  `json:"eventTypeOther,omitempty" validate:"omitempty,max=100"`
	EventName           *string  `json:"eventName,omitempty" validate:"omitempty,max=200"`
	EventDate           string   `json:"eventDate" validate:"required"`
	EventTime           string   `json:"eventTime" validate:"required"`
	DurationHours       float64  `json:"durationHours" validate:"required,gt=0,lte=24"`
	GuestCount          int      `json:"guestCount" validate:"required,min=1,max=1000"`
	VenueName           *string  `json:"venueName,omitempty" validate:"omitempty,max=200"`
	VenueAddress        string   `json:"venueAddress" validate:"required,min=5,max=500"`
	VenueArea           string   `json:"venueArea" validate:"required,min=2,max=100"`
	VenueState          string   `json:"venueState" validate:"required,min=2,max=100"`
	VenueNotes          *string  `json:"venueNotes,omitempty" validate:"omitempty,max=500"`
	ServiceStyle        string   `json:"serviceStyle,omitempty" validate:"omitempty,oneof=drop_off buffet full_service"`
	SelectedItemIDs     []int    `json:"selectedItemIds" validate:"max=20,dive,min=1"`
	MenuNotes           *string  `json:"menuNotes,omitempty" validate:"omitempty,max=1000"`
	DietaryRequirements []string `json:"dietaryRequirements" validate:"max=7,dive,oneof=vegetarian vegan halal no_pork nut_free gluten_free dairy_free"`
	DietaryNotes        *string  `json:"dietaryNotes,omitempty" validate:"omitempty,max=500"`
	Budget              *float64 `json:"budget,omitempty" validate:"omitempty,min=0"`
	SpecialRequests     *string  `json:"specialRequests,omitempty" validate:"omitempty,max=1000"`
}

// CateringLookupRequest lets a customer find their request.
type CateringLookupRequest struct {
	Reference string `json:"reference" validate:"required,max=30"`
	Phone     string `json:"phone" validate:"required,max=30"`
}

// CateringRequestUpdate is the admin's follow-up on a request; nil fields are
// left unchanged.
type CateringRequestUpdate struct {
	Status       *string  `json:"status,omitempty"`
	QuotedAmount *float64 `json:"quotedAmount,omitempty" validate:"omitempty,min=0"`
	AdminNotes   *string  `json:"adminNotes,omitempty" validate:"omitempty,max=2000"`
}
