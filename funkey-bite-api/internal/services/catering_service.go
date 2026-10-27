package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// DefaultCustomNoticeDays is the minimum notice for custom requests, which
// have no package to set one.
const DefaultCustomNoticeDays = 3

// maxCateringDaysAhead limits how far ahead an event can be booked.
const maxCateringDaysAhead = 365

var (
	ErrCateringRequestNotFound = errors.New("catering request not found")
	ErrCateringPackageNotFound = errors.New("catering package not found")
)

// CateringValidationError explains why a request or package was rejected;
// its message is safe to show to the person who submitted it.
type CateringValidationError struct {
	Message string
}

func (e *CateringValidationError) Error() string { return e.Message }

func invalid(format string, args ...any) error {
	return &CateringValidationError{Message: fmt.Sprintf(format, args...)}
}

type CateringService interface {
	CreateRequest(input models.CateringRequestInput, userID *int) (*models.CateringRequest, error)
	GetRequestByID(id int) (*models.CateringRequest, error)
	LookupRequest(reference, phone string) (*models.CateringRequest, error)
	GetUserRequests(userID int) ([]models.CateringRequest, error)
	GetAllRequests() ([]models.CateringRequest, error)
	UpdateRequestStatus(id int, status string) error
	UpdateRequest(id int, update *models.CateringRequestUpdate) (*models.CateringRequest, error)

	ListPackages(activeOnly bool) ([]models.CateringPackage, error)
	GetPackage(id int) (*models.CateringPackage, error)
	CreatePackage(input *models.CateringPackageInput) (*models.CateringPackage, error)
	UpdatePackage(id int, input *models.CateringPackageInput) (*models.CateringPackage, error)
}

type cateringRequestStore interface {
	Create(request *models.CateringRequest) (*models.CateringRequest, error)
	GetByID(id int) (*models.CateringRequest, error)
	GetByReference(reference string) (*models.CateringRequest, error)
	GetByUserID(userID int) ([]models.CateringRequest, error)
	GetAll() ([]models.CateringRequest, error)
	UpdateStatus(id int, status string) error
	ApplyUpdate(id int, update *models.CateringRequestUpdate) error
}

type cateringPackageStore interface {
	List(activeOnly bool) ([]models.CateringPackage, error)
	GetByID(id int) (*models.CateringPackage, error)
	Create(input *models.CateringPackageInput) (*models.CateringPackage, error)
	Update(id int, input *models.CateringPackageInput) (*models.CateringPackage, error)
}

type cateringService struct {
	cateringRepo        cateringRequestStore
	packageRepo         cateringPackageStore
	menuRepo            menuItemReader
	notificationService NotificationService
	location            *time.Location
	now                 func() time.Time
}

func NewCateringService(
	cateringRepo cateringRequestStore,
	packageRepo cateringPackageStore,
	menuRepo menuItemReader,
	notificationService NotificationService,
) CateringService {
	return &cateringService{
		cateringRepo:        cateringRepo,
		packageRepo:         packageRepo,
		menuRepo:            menuRepo,
		notificationService: notificationService,
		location:            LoadBusinessLocation(),
		now:                 time.Now,
	}
}

func (s *cateringService) CreateRequest(input models.CateringRequestInput, userID *int) (*models.CateringRequest, error) {
	phone := utils.CleanPhone(input.ContactPhone)
	if !utils.IsValidPhone(phone) {
		return nil, invalid("Enter a valid phone number")
	}

	eventTypeOther := trimmedPtr(input.EventTypeOther)
	if input.EventType == "other" && (eventTypeOther == nil || len(*eventTypeOther) < 2) {
		return nil, invalid("Tell us what kind of event it is")
	}
	if input.EventType != "other" {
		eventTypeOther = nil
	}

	eventDate, err := time.ParseInLocation("2006-01-02", input.EventDate, s.location)
	if err != nil {
		return nil, invalid("Event date must be in YYYY-MM-DD format")
	}
	if _, ok := clockMinutes(input.EventTime); !ok {
		return nil, invalid("Event time must be in 24-hour HH:MM format")
	}

	request := &models.CateringRequest{
		RequestType:         input.RequestType,
		UserID:              userID,
		ContactName:         strings.TrimSpace(input.ContactName),
		ContactPhone:        phone,
		ContactEmail:        trimmedPtr(&input.ContactEmail),
		Organization:        trimmedPtr(input.Organization),
		EventType:           input.EventType,
		EventTypeOther:      eventTypeOther,
		EventName:           trimmedPtr(input.EventName),
		EventDate:           input.EventDate,
		EventTime:           &input.EventTime,
		DurationHours:       &input.DurationHours,
		GuestCount:          input.GuestCount,
		VenueName:           trimmedPtr(input.VenueName),
		VenueAddress:        trimmedPtr(&input.VenueAddress),
		VenueArea:           trimmedPtr(&input.VenueArea),
		VenueState:          trimmedPtr(&input.VenueState),
		VenueNotes:          trimmedPtr(input.VenueNotes),
		MenuNotes:           trimmedPtr(input.MenuNotes),
		DietaryRequirements: uniqueStrings(input.DietaryRequirements),
		DietaryNotes:        trimmedPtr(input.DietaryNotes),
		Budget:              input.Budget,
		SpecialRequests:     trimmedPtr(input.SpecialRequests),
		Status:              models.CateringStatusPending,
		CreatedAt:           s.now(),
	}

	itemIDs := uniqueInts(input.SelectedItemIDs)
	var noticeDays int

	switch input.RequestType {
	case models.CateringRequestTypePackage:
		pkg, err := s.packageForRequest(input.PackageID)
		if err != nil {
			return nil, err
		}
		if input.GuestCount < pkg.MinGuests {
			return nil, invalid("The %s package is for at least %d guests", pkg.Name, pkg.MinGuests)
		}
		if pkg.MaxGuests != nil && input.GuestCount > *pkg.MaxGuests {
			return nil, invalid("The %s package is for at most %d guests", pkg.Name, *pkg.MaxGuests)
		}
		selected, err := packageSelection(pkg, itemIDs)
		if err != nil {
			return nil, err
		}

		estimate := pkg.PricePerPerson * float64(input.GuestCount)
		request.PackageID = &pkg.ID
		request.PackageName = &pkg.Name
		request.ServiceStyle = &pkg.ServiceStyle
		request.SelectedItems = selected
		request.EstimatedTotal = &estimate
		noticeDays = pkg.MinNoticeDays

	case models.CateringRequestTypeCustom:
		if input.ServiceStyle == "" {
			return nil, invalid("Choose how you would like the food served")
		}
		selected, err := s.menuSelection(itemIDs)
		if err != nil {
			return nil, err
		}
		if len(selected) == 0 && (request.MenuNotes == nil || len(*request.MenuNotes) < 10) {
			return nil, invalid("Tell us what food you would like, or pick dishes from the menu")
		}

		style := input.ServiceStyle
		request.ServiceStyle = &style
		request.SelectedItems = selected
		noticeDays = DefaultCustomNoticeDays
	}

	if err := s.checkEventDate(eventDate, noticeDays); err != nil {
		return nil, err
	}

	created, err := s.createWithReference(request)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.notificationService.SendCateringConfirmation(created); err != nil {
			log.Printf("Failed to send catering confirmation for %s: %v", created.Reference, err)
		}
	}()

	return created, nil
}

func (s *cateringService) packageForRequest(packageID *int) (*models.CateringPackage, error) {
	if packageID == nil {
		return nil, invalid("Choose a catering package")
	}
	pkg, err := s.packageRepo.GetByID(*packageID)
	if err != nil {
		return nil, err
	}
	if pkg == nil || !pkg.IsActive {
		return nil, invalid("That catering package is no longer available")
	}
	return pkg, nil
}

// packageSelection checks the dishes chosen for a package: exactly the number
// it offers, all from its list and currently available.
func packageSelection(pkg *models.CateringPackage, itemIDs []int) ([]models.CateringSelectedItem, error) {
	if pkg.DishChoices == 0 {
		return []models.CateringSelectedItem{}, nil
	}
	if len(itemIDs) != pkg.DishChoices {
		return nil, invalid("Choose %d dishes for the %s package", pkg.DishChoices, pkg.Name)
	}

	offered := make(map[int]models.CateringPackageItem, len(pkg.MenuItems))
	for _, item := range pkg.MenuItems {
		offered[item.MenuItemID] = item
	}

	selected := make([]models.CateringSelectedItem, 0, len(itemIDs))
	for _, id := range itemIDs {
		item, ok := offered[id]
		if !ok {
			return nil, invalid("One of the chosen dishes is not part of the %s package", pkg.Name)
		}
		if !item.IsAvailable {
			return nil, invalid("%s is not available right now; please choose another dish", item.Name)
		}
		selected = append(selected, models.CateringSelectedItem{MenuItemID: id, Name: item.Name})
	}
	return selected, nil
}

// menuSelection checks dish preferences on a custom request.
func (s *cateringService) menuSelection(itemIDs []int) ([]models.CateringSelectedItem, error) {
	selected := make([]models.CateringSelectedItem, 0, len(itemIDs))
	for _, id := range itemIDs {
		item, err := s.menuRepo.GetByID(id)
		if err != nil {
			return nil, err
		}
		if item == nil || !item.IsAvailable {
			return nil, invalid("One of the chosen dishes is not available; please choose another")
		}
		selected = append(selected, models.CateringSelectedItem{MenuItemID: id, Name: item.Name})
	}
	return selected, nil
}

// checkEventDate enforces the notice period, counted in whole days in the
// business timezone.
func (s *cateringService) checkEventDate(eventDate time.Time, noticeDays int) error {
	now := s.now().In(s.location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.location)
	earliest := today.AddDate(0, 0, noticeDays)

	if eventDate.Before(today) {
		return invalid("The event date has already passed")
	}
	if eventDate.Before(earliest) {
		return invalid("We need at least %d days' notice; the earliest date we can cater is %s",
			noticeDays, earliest.Format("2 Jan 2006"))
	}
	if eventDate.After(today.AddDate(0, 0, maxCateringDaysAhead)) {
		return invalid("Events can be booked up to a year ahead")
	}
	return nil
}

func (s *cateringService) createWithReference(request *models.CateringRequest) (*models.CateringRequest, error) {
	for attempt := 0; attempt < 3; attempt++ {
		reference, err := newCateringReference()
		if err != nil {
			return nil, err
		}
		request.Reference = reference

		created, err := s.cateringRepo.Create(request)
		if errors.Is(err, repository.ErrDuplicateCateringReference) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to create catering request: %w", err)
		}
		return created, nil
	}
	return nil, fmt.Errorf("failed to create catering request: could not generate a unique reference")
}

// referenceAlphabet leaves out characters that are easy to misread (0/O, 1/I).
const referenceAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// newCateringReference returns a short random reference such as "CAT-7K3M9Q".
func newCateringReference() (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating catering reference: %w", err)
	}
	for i, b := range buf {
		buf[i] = referenceAlphabet[int(b)%len(referenceAlphabet)]
	}
	return "CAT-" + string(buf), nil
}

func (s *cateringService) GetRequestByID(id int) (*models.CateringRequest, error) {
	request, err := s.cateringRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get catering request: %w", err)
	}
	if request == nil {
		return nil, ErrCateringRequestNotFound
	}
	return request, nil
}

// LookupRequest finds a request for the customer who made it: the reference
// and the phone number must both match.
func (s *cateringService) LookupRequest(reference, phone string) (*models.CateringRequest, error) {
	request, err := s.cateringRepo.GetByReference(reference)
	if err != nil {
		return nil, fmt.Errorf("failed to get catering request: %w", err)
	}
	if request == nil || !utils.SamePhone(request.ContactPhone, phone) {
		return nil, ErrCateringRequestNotFound
	}
	return request, nil
}

func (s *cateringService) GetUserRequests(userID int) ([]models.CateringRequest, error) {
	requests, err := s.cateringRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user catering requests: %w", err)
	}
	return requests, nil
}

func (s *cateringService) GetAllRequests() ([]models.CateringRequest, error) {
	requests, err := s.cateringRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all catering requests: %w", err)
	}
	return requests, nil
}

func (s *cateringService) UpdateRequestStatus(id int, status string) error {
	_, err := s.UpdateRequest(id, &models.CateringRequestUpdate{Status: &status})
	return err
}

// UpdateRequest applies the admin's follow-up. A request can only be marked
// quoted once it has a quoted amount.
func (s *cateringService) UpdateRequest(id int, update *models.CateringRequestUpdate) (*models.CateringRequest, error) {
	if update.Status != nil && !models.IsValidCateringStatus(*update.Status) {
		return nil, invalid("Unknown status %q", *update.Status)
	}

	current, err := s.GetRequestByID(id)
	if err != nil {
		return nil, err
	}

	if update.Status != nil && *update.Status == string(models.CateringStatusQuoted) &&
		update.QuotedAmount == nil && current.QuotedAmount == nil {
		return nil, invalid("Enter the quoted amount before marking the request as quoted")
	}
	if update.AdminNotes != nil {
		notes := strings.TrimSpace(*update.AdminNotes)
		update.AdminNotes = &notes
	}

	if err := s.cateringRepo.ApplyUpdate(id, update); err != nil {
		return nil, fmt.Errorf("failed to update catering request: %w", err)
	}
	return s.GetRequestByID(id)
}

func (s *cateringService) ListPackages(activeOnly bool) ([]models.CateringPackage, error) {
	return s.packageRepo.List(activeOnly)
}

func (s *cateringService) GetPackage(id int) (*models.CateringPackage, error) {
	pkg, err := s.packageRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, ErrCateringPackageNotFound
	}
	return pkg, nil
}

func (s *cateringService) CreatePackage(input *models.CateringPackageInput) (*models.CateringPackage, error) {
	if err := s.normalizePackageInput(input); err != nil {
		return nil, err
	}
	return s.packageRepo.Create(input)
}

func (s *cateringService) UpdatePackage(id int, input *models.CateringPackageInput) (*models.CateringPackage, error) {
	if err := s.normalizePackageInput(input); err != nil {
		return nil, err
	}
	pkg, err := s.packageRepo.Update(id, input)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, ErrCateringPackageNotFound
	}
	return pkg, nil
}

// normalizePackageInput trims the package's text and checks the rules that
// span fields: the guest range, and that customers can choose enough dishes.
func (s *cateringService) normalizePackageInput(input *models.CateringPackageInput) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	inclusions := make([]string, 0, len(input.Inclusions))
	for _, inclusion := range input.Inclusions {
		if trimmed := strings.TrimSpace(inclusion); trimmed != "" {
			inclusions = append(inclusions, trimmed)
		}
	}
	input.Inclusions = inclusions

	if input.MaxGuests != nil && *input.MaxGuests < input.MinGuests {
		return invalid("Maximum guests cannot be less than minimum guests")
	}

	input.MenuItemIDs = uniqueInts(input.MenuItemIDs)
	for _, id := range input.MenuItemIDs {
		item, err := s.menuRepo.GetByID(id)
		if err != nil {
			return err
		}
		if item == nil {
			return invalid("Menu item %d does not exist", id)
		}
	}
	if input.DishChoices > len(input.MenuItemIDs) {
		return invalid("Customers pick %d dishes, so link at least %d menu items", input.DishChoices, input.DishChoices)
	}
	return nil
}

func trimmedPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func uniqueInts(values []int) []int {
	seen := make(map[int]bool, len(values))
	out := make([]int, 0, len(values))
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
