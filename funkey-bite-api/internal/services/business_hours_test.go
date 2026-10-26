package services

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

func week(open, close string) []models.OpeningHours {
	hours := make([]models.OpeningHours, 0, 7)
	for _, day := range weekdayNames {
		hours = append(hours, models.OpeningHours{Day: day, OpenTime: open, CloseTime: close, IsOpen: true})
	}
	return hours
}

func TestLoadBusinessLocation(t *testing.T) {
	t.Setenv("BUSINESS_TIMEZONE", "")
	if got := LoadBusinessLocation().String(); got != DefaultBusinessTimezone {
		t.Fatalf("default location = %q, want %q", got, DefaultBusinessTimezone)
	}

	t.Setenv("BUSINESS_TIMEZONE", "Not/AZone")
	if got := LoadBusinessLocation().String(); got != DefaultBusinessTimezone {
		t.Fatalf("invalid zone must fall back, got %q", got)
	}
}

func TestIsOpenAtUsesBusinessTimezone(t *testing.T) {
	lagos := LoadBusinessLocation()
	hours := week("08:00", "22:00")

	// 21:30 UTC is 22:30 in Lagos: a UTC server comparing its own clock would
	// wrongly report the business as open.
	utc := time.Date(2026, 10, 26, 21, 30, 0, 0, time.UTC)
	if isOpenAt(hours, utc.In(lagos)) {
		t.Fatal("expected closed at 22:30 Lagos time")
	}
	if !isOpenAt(hours, utc.Add(-time.Hour).In(lagos)) {
		t.Fatal("expected open at 21:30 Lagos time")
	}
}

func TestIsOpenAt(t *testing.T) {
	// 2026-10-26 is a Monday.
	at := func(day int, clock string) time.Time {
		parsed, _ := time.Parse("15:04", clock)
		return time.Date(2026, 10, 26+day, parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
	}

	overnightFriday := week("08:00", "22:00")
	overnightFriday[5] = models.OpeningHours{Day: "Friday", OpenTime: "18:00", CloseTime: "02:00", IsOpen: true}
	closedSunday := week("08:00", "22:00")
	closedSunday[0].IsOpen = false
	untilMidnight := week("08:00", "00:00")

	cases := []struct {
		name  string
		hours []models.OpeningHours
		at    time.Time
		want  bool
	}{
		{"at opening time", week("08:00", "22:00"), at(0, "08:00"), true},
		{"before opening", week("08:00", "22:00"), at(0, "07:59"), false},
		{"closing time is exclusive", week("08:00", "22:00"), at(0, "22:00"), false},
		{"overnight evening", overnightFriday, at(4, "23:30"), true},
		{"overnight spills into next morning", overnightFriday, at(5, "01:59"), true},
		{"overnight ends at close", overnightFriday, at(5, "02:00"), false},
		{"normal day does not spill over", overnightFriday, at(4, "01:00"), false},
		{"closed day", closedSunday, at(6, "12:00"), false},
		{"open until midnight", untilMidnight, at(0, "23:59"), true},
		{"midnight close ends the day", untilMidnight, at(1, "00:00"), false},
		{"missing day is closed", week("08:00", "22:00")[1:], at(6, "12:00"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isOpenAt(tc.hours, tc.at); got != tc.want {
				t.Fatalf("isOpenAt(%s %s) = %t, want %t", tc.at.Weekday(), tc.at.Format("15:04"), got, tc.want)
			}
		})
	}
}

func TestNormalizeOpeningHours(t *testing.T) {
	// Lower-case names in any order come back canonical and Monday first.
	input := week("09:00", "21:00")
	input[0].Day = "sunday"
	input[0], input[3] = input[3], input[0]
	input[1] = models.OpeningHours{Day: "Monday", IsOpen: false}

	got, err := normalizeOpeningHours(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var days []string
	for _, h := range got {
		days = append(days, h.Day)
	}
	if strings.Join(days, ",") != "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday" {
		t.Fatalf("unexpected day order %v", days)
	}

	cases := []struct {
		name    string
		mutate  func([]models.OpeningHours) []models.OpeningHours
		wantErr string
	}{
		{"missing day", func(h []models.OpeningHours) []models.OpeningHours { return h[1:] }, "all 7 days"},
		{"duplicate day", func(h []models.OpeningHours) []models.OpeningHours { h[1].Day = "Sunday"; return h }, "more than once"},
		{"unknown day", func(h []models.OpeningHours) []models.OpeningHours { h[1].Day = "Funday"; return h }, "unknown day"},
		{"bad time", func(h []models.OpeningHours) []models.OpeningHours { h[1].OpenTime = "9am"; return h }, "HH:MM"},
		{"out of range time", func(h []models.OpeningHours) []models.OpeningHours { h[1].CloseTime = "24:00"; return h }, "HH:MM"},
		{"same open and close", func(h []models.OpeningHours) []models.OpeningHours { h[1].CloseTime = h[1].OpenTime; return h }, "cannot be the same"},
		{"bad time on closed day", func(h []models.OpeningHours) []models.OpeningHours {
			h[1].IsOpen, h[1].OpenTime = false, "9am"
			return h
		}, "HH:MM"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := normalizeOpeningHours(tc.mutate(week("08:00", "22:00")))
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestSettingsUpdateAcceptsAdminPayload(t *testing.T) {
	// The admin app sends openingHours as an array; it used to be decoded into
	// a string, which rejected every save.
	body := `{
		"businessName": "  Funkey Grab & Bite ",
		"phoneNumber": "+234 801 234 5678",
		"email": "hello@example.com",
		"address": "12 Admiralty Way",
		"city": "Lekki",
		"state": "Lagos",
		"postalCode": "",
		"country": "Nigeria",
		"openingHours": [
			{"day": "Monday", "openTime": "08:00", "closeTime": "22:00", "isOpen": true},
			{"day": "Tuesday", "openTime": "08:00", "closeTime": "22:00", "isOpen": true},
			{"day": "Wednesday", "openTime": "08:00", "closeTime": "22:00", "isOpen": true},
			{"day": "Thursday", "openTime": "08:00", "closeTime": "22:00", "isOpen": true},
			{"day": "Friday", "openTime": "18:00", "closeTime": "02:00", "isOpen": true},
			{"day": "Saturday", "openTime": "09:00", "closeTime": "23:00", "isOpen": true},
			{"day": "Sunday", "openTime": "", "closeTime": "", "isOpen": false}
		]
	}`

	var update models.BusinessSettingsUpdate
	if err := json.Unmarshal([]byte(body), &update); err != nil {
		t.Fatalf("admin payload must decode: %v", err)
	}
	if err := normalizeSettingsUpdate(&update); err != nil {
		t.Fatalf("admin payload must validate: %v", err)
	}
	if *update.BusinessName != "Funkey Grab & Bite" {
		t.Fatalf("expected trimmed business name, got %q", *update.BusinessName)
	}
	if len(*update.OpeningHours) != 7 {
		t.Fatalf("expected 7 days, got %d", len(*update.OpeningHours))
	}
}

func TestNormalizeSettingsUpdateRejections(t *testing.T) {
	str := func(s string) *string { return &s }
	cases := []struct {
		name    string
		update  models.BusinessSettingsUpdate
		wantErr string
	}{
		{"blank business name", models.BusinessSettingsUpdate{BusinessName: str("   ")}, "business name is required"},
		{"blank street", models.BusinessSettingsUpdate{Address: str("")}, "street address is required"},
		{"blank city", models.BusinessSettingsUpdate{City: str("")}, "city is required"},
		{"long state", models.BusinessSettingsUpdate{State: str(strings.Repeat("a", 101))}, "at most 100"},
		{"bad phone", models.BusinessSettingsUpdate{PhoneNumber: str("call me")}, "phone number"},
		{"bad email", models.BusinessSettingsUpdate{Email: str("not-an-email")}, "valid email"},
		{"email with display name", models.BusinessSettingsUpdate{Email: str("Shop <shop@example.com>")}, "valid email"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := normalizeSettingsUpdate(&tc.update)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}

	if err := normalizeSettingsUpdate(&models.BusinessSettingsUpdate{PostalCode: str("")}); err != nil {
		t.Fatalf("postal code is optional, got %v", err)
	}
}
