package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type memoryCateringStore struct {
	requests map[int]*models.CateringRequest
	nextID   int
}

func newMemoryCateringStore() *memoryCateringStore {
	return &memoryCateringStore{requests: map[int]*models.CateringRequest{}, nextID: 1}
}

func (m *memoryCateringStore) Create(request *models.CateringRequest) (*models.CateringRequest, error) {
	request.ID = m.nextID
	m.nextID++
	copy := *request
	m.requests[request.ID] = &copy
	return request, nil
}
func (m *memoryCateringStore) GetByID(id int) (*models.CateringRequest, error) {
	if r, ok := m.requests[id]; ok {
		copy := *r
		return &copy, nil
	}
	return nil, nil
}
func (m *memoryCateringStore) GetByReference(reference string) (*models.CateringRequest, error) {
	for _, r := range m.requests {
		if r.Reference == reference {
			copy := *r
			return &copy, nil
		}
	}
	return nil, nil
}
func (m *memoryCateringStore) GetByUserID(int) ([]models.CateringRequest, error) { return nil, nil }
func (m *memoryCateringStore) GetAll() ([]models.CateringRequest, error)         { return nil, nil }
func (m *memoryCateringStore) UpdateStatus(id int, status string) error {
	return m.ApplyUpdate(id, &models.CateringRequestUpdate{Status: &status})
}
func (m *memoryCateringStore) ApplyUpdate(id int, update *models.CateringRequestUpdate) error {
	r, ok := m.requests[id]
	if !ok {
		return errors.New("not found")
	}
	if update.Status != nil {
		r.Status = models.CateringStatus(*update.Status)
	}
	if update.QuotedAmount != nil {
		r.QuotedAmount = update.QuotedAmount
	}
	if update.AdminNotes != nil {
		r.AdminNotes = update.AdminNotes
	}
	return nil
}

type memoryPackageStore struct {
	packages map[int]*models.CateringPackage
	saved    *models.CateringPackageInput
}

func (m *memoryPackageStore) List(bool) ([]models.CateringPackage, error) { return nil, nil }
func (m *memoryPackageStore) GetByID(id int) (*models.CateringPackage, error) {
	if p, ok := m.packages[id]; ok {
		copy := *p
		return &copy, nil
	}
	return nil, nil
}
func (m *memoryPackageStore) Create(input *models.CateringPackageInput) (*models.CateringPackage, error) {
	m.saved = input
	return &models.CateringPackage{ID: 99, Name: input.Name}, nil
}
func (m *memoryPackageStore) Update(id int, input *models.CateringPackageInput) (*models.CateringPackage, error) {
	m.saved = input
	if _, ok := m.packages[id]; !ok {
		return nil, nil
	}
	return &models.CateringPackage{ID: id, Name: input.Name}, nil
}

type memoryMenu map[int]*models.MenuItem

func (m memoryMenu) GetByID(id int) (*models.MenuItem, error) { return m[id], nil }

type silentNotifier struct{ NotificationService }

func (silentNotifier) SendCateringConfirmation(*models.CateringRequest) error { return nil }

// 2026-10-27 10:00 in Lagos.
var cateringNow = time.Date(2026, 10, 27, 9, 0, 0, 0, time.UTC)

func newTestCateringService() (*cateringService, *memoryCateringStore, *memoryPackageStore) {
	maxGuests := 150
	packages := &memoryPackageStore{packages: map[int]*models.CateringPackage{
		1: {
			ID: 1, Name: "Owambe", PricePerPerson: 5000, MinGuests: 50, MaxGuests: &maxGuests,
			ServiceStyle: models.ServiceStyleBuffet, MinNoticeDays: 7, DishChoices: 2, IsActive: true,
			MenuItems: []models.CateringPackageItem{
				{MenuItemID: 10, Name: "Jollof Rice", IsAvailable: true},
				{MenuItemID: 11, Name: "Fried Rice", IsAvailable: true},
				{MenuItemID: 12, Name: "Asun", IsAvailable: false},
			},
		},
		2: {ID: 2, Name: "Retired", MinGuests: 1, ServiceStyle: models.ServiceStyleDropOff, IsActive: false},
	}}
	menu := memoryMenu{
		10: {ID: 10, Name: "Jollof Rice", IsAvailable: true},
		20: {ID: 20, Name: "Small Chops", IsAvailable: true},
		21: {ID: 21, Name: "Moi Moi", IsAvailable: false},
	}
	store := newMemoryCateringStore()
	svc := NewCateringService(store, packages, menu, silentNotifier{}).(*cateringService)
	svc.now = func() time.Time { return cateringNow }
	return svc, store, packages
}

func packageInput() models.CateringRequestInput {
	packageID := 1
	return models.CateringRequestInput{
		RequestType: models.CateringRequestTypePackage, PackageID: &packageID,
		ContactName: "Ada Obi", ContactPhone: "0801 234 5678", ContactEmail: " ada@example.com ",
		EventType: "party", EventDate: "2026-11-10", EventTime: "14:00", DurationHours: 5,
		GuestCount: 80, VenueAddress: "5 Bourdillon Road", VenueArea: "Ikoyi", VenueState: "Lagos",
		SelectedItemIDs: []int{10, 11}, DietaryRequirements: []string{"halal", "halal"},
	}
}

func customInput() models.CateringRequestInput {
	input := packageInput()
	notes := "Small chops, peppered fish and fruit platters"
	input.RequestType = models.CateringRequestTypeCustom
	input.PackageID = nil
	input.SelectedItemIDs = nil
	input.ServiceStyle = models.ServiceStyleFullService
	input.MenuNotes = &notes
	input.EventDate = "2026-10-31"
	return input
}

func TestCreatePackageRequest(t *testing.T) {
	svc, _, _ := newTestCateringService()

	request, err := svc.CreateRequest(packageInput(), nil)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}
	if !strings.HasPrefix(request.Reference, "CAT-") || len(request.Reference) != 10 {
		t.Fatalf("unexpected reference %q", request.Reference)
	}
	if request.ContactPhone != "08012345678" || *request.ContactEmail != "ada@example.com" {
		t.Fatalf("contact details not cleaned: %q %q", request.ContactPhone, *request.ContactEmail)
	}
	if *request.ServiceStyle != models.ServiceStyleBuffet || *request.EstimatedTotal != 400000 {
		t.Fatalf("expected package service style and 80 x 5000 estimate, got %v %v", *request.ServiceStyle, *request.EstimatedTotal)
	}
	if len(request.SelectedItems) != 2 || request.SelectedItems[0].Name != "Jollof Rice" {
		t.Fatalf("unexpected selected items %+v", request.SelectedItems)
	}
	if len(request.DietaryRequirements) != 1 {
		t.Fatalf("expected dietary requirements de-duplicated, got %v", request.DietaryRequirements)
	}
}

func TestCreatePackageRequestRejections(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*models.CateringRequestInput)
		wantErr string
	}{
		{"no package", func(in *models.CateringRequestInput) { in.PackageID = nil }, "Choose a catering package"},
		{"inactive package", func(in *models.CateringRequestInput) { id := 2; in.PackageID = &id }, "no longer available"},
		{"too few guests", func(in *models.CateringRequestInput) { in.GuestCount = 20 }, "at least 50 guests"},
		{"too many guests", func(in *models.CateringRequestInput) { in.GuestCount = 200 }, "at most 150 guests"},
		{"wrong dish count", func(in *models.CateringRequestInput) { in.SelectedItemIDs = []int{10} }, "Choose 2 dishes"},
		{"dish not in package", func(in *models.CateringRequestInput) { in.SelectedItemIDs = []int{10, 20} }, "not part of the Owambe package"},
		{"unavailable dish", func(in *models.CateringRequestInput) { in.SelectedItemIDs = []int{10, 12} }, "Asun is not available"},
		{"short notice", func(in *models.CateringRequestInput) { in.EventDate = "2026-10-30" }, "at least 7 days' notice; the earliest date we can cater is 3 Nov 2026"},
		{"past date", func(in *models.CateringRequestInput) { in.EventDate = "2026-10-20" }, "already passed"},
		{"too far ahead", func(in *models.CateringRequestInput) { in.EventDate = "2028-01-01" }, "up to a year ahead"},
		{"bad date", func(in *models.CateringRequestInput) { in.EventDate = "10/11/2026" }, "YYYY-MM-DD"},
		{"bad time", func(in *models.CateringRequestInput) { in.EventTime = "2pm" }, "HH:MM"},
		{"bad phone", func(in *models.CateringRequestInput) { in.ContactPhone = "12-34" }, "valid phone number"},
		{"other event without description", func(in *models.CateringRequestInput) { in.EventType = "other" }, "what kind of event"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, store, _ := newTestCateringService()
			input := packageInput()
			tc.mutate(&input)

			_, err := svc.CreateRequest(input, nil)
			var invalid *CateringValidationError
			if !errors.As(err, &invalid) || !strings.Contains(invalid.Message, tc.wantErr) {
				t.Fatalf("expected validation error containing %q, got %v", tc.wantErr, err)
			}
			if len(store.requests) != 0 {
				t.Fatal("rejected request must not be saved")
			}
		})
	}
}

func TestCreateCustomRequest(t *testing.T) {
	svc, _, _ := newTestCateringService()

	input := customInput()
	input.SelectedItemIDs = []int{20}
	request, err := svc.CreateRequest(input, nil)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}
	if request.PackageID != nil || request.EstimatedTotal != nil {
		t.Fatalf("custom requests have no package or estimate, got %+v", request)
	}
	if *request.ServiceStyle != models.ServiceStyleFullService || request.SelectedItems[0].Name != "Small Chops" {
		t.Fatalf("unexpected custom request %+v", request)
	}

	cases := []struct {
		name    string
		mutate  func(*models.CateringRequestInput)
		wantErr string
	}{
		{"no service style", func(in *models.CateringRequestInput) { in.ServiceStyle = "" }, "how you would like the food served"},
		{"no menu at all", func(in *models.CateringRequestInput) { in.MenuNotes = nil }, "what food you would like"},
		{"unavailable preference", func(in *models.CateringRequestInput) { in.SelectedItemIDs = []int{21} }, "not available"},
		{"inside default notice", func(in *models.CateringRequestInput) { in.EventDate = "2026-10-29" }, "at least 3 days' notice"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, _, _ := newTestCateringService()
			input := customInput()
			tc.mutate(&input)
			_, err := svc.CreateRequest(input, nil)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestLookupRequestNeedsMatchingPhone(t *testing.T) {
	svc, _, _ := newTestCateringService()
	request, err := svc.CreateRequest(packageInput(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if found, err := svc.LookupRequest(request.Reference, "+234 801 234 5678"); err != nil || found.ID != request.ID {
		t.Fatalf("lookup with the same phone in another format should succeed, got %v", err)
	}
	if _, err := svc.LookupRequest(request.Reference, "08099999999"); !errors.Is(err, ErrCateringRequestNotFound) {
		t.Fatalf("wrong phone must look like a missing request, got %v", err)
	}
}

func TestUpdateRequestQuoteRules(t *testing.T) {
	svc, _, _ := newTestCateringService()
	request, err := svc.CreateRequest(packageInput(), nil)
	if err != nil {
		t.Fatal(err)
	}

	quoted := string(models.CateringStatusQuoted)
	if _, err := svc.UpdateRequest(request.ID, &models.CateringRequestUpdate{Status: &quoted}); err == nil {
		t.Fatal("quoting without an amount must be rejected")
	}

	amount := 450000.0
	notes := "  Includes two servers  "
	updated, err := svc.UpdateRequest(request.ID, &models.CateringRequestUpdate{Status: &quoted, QuotedAmount: &amount, AdminNotes: &notes})
	if err != nil {
		t.Fatalf("UpdateRequest() error = %v", err)
	}
	if updated.Status != models.CateringStatusQuoted || *updated.QuotedAmount != 450000 || *updated.AdminNotes != "Includes two servers" {
		t.Fatalf("unexpected update result %+v", updated)
	}

	bogus := "shipped"
	if _, err := svc.UpdateRequest(request.ID, &models.CateringRequestUpdate{Status: &bogus}); err == nil {
		t.Fatal("unknown status must be rejected")
	}
	if _, err := svc.UpdateRequest(999, &models.CateringRequestUpdate{Status: &quoted}); !errors.Is(err, ErrCateringRequestNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPackageInputRules(t *testing.T) {
	svc, _, packages := newTestCateringService()
	maxGuests := 10

	cases := []struct {
		name    string
		input   models.CateringPackageInput
		wantErr string
	}{
		{"max below min", models.CateringPackageInput{Name: "A", MinGuests: 20, MaxGuests: &maxGuests}, "cannot be less than minimum"},
		{"unknown menu item", models.CateringPackageInput{Name: "A", MinGuests: 1, MenuItemIDs: []int{404}}, "does not exist"},
		{"too few dishes to choose from", models.CateringPackageInput{Name: "A", MinGuests: 1, DishChoices: 2, MenuItemIDs: []int{10}}, "link at least 2"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.input
			if _, err := svc.CreatePackage(&input); err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}

	input := models.CateringPackageInput{
		Name: "  Office Lunch ", MinGuests: 10, DishChoices: 1, MenuItemIDs: []int{10, 10, 20},
		Inclusions: []string{" Bottled water ", "", "Disposable packs"},
	}
	if _, err := svc.CreatePackage(&input); err != nil {
		t.Fatalf("CreatePackage() error = %v", err)
	}
	if packages.saved.Name != "Office Lunch" || len(packages.saved.MenuItemIDs) != 2 || len(packages.saved.Inclusions) != 2 {
		t.Fatalf("package input not normalized: %+v", packages.saved)
	}

	if _, err := svc.UpdatePackage(404, &input); !errors.Is(err, ErrCateringPackageNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
