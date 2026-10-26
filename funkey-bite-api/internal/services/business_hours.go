package services

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	// Embeds the IANA timezone database so Africa/Lagos resolves in minimal
	// container images (alpine) that ship without /usr/share/zoneinfo.
	_ "time/tzdata"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// DefaultBusinessTimezone is used when BUSINESS_TIMEZONE is unset or invalid.
const DefaultBusinessTimezone = "Africa/Lagos"

var (
	clockPattern = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)
	phonePattern = regexp.MustCompile(`^\+?[0-9][0-9 ()-]{5,18}$`)
)

// weekdayNames are indexed by time.Weekday (Sunday = 0).
var weekdayNames = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// LoadBusinessLocation returns the timezone the opening hours are written in.
// Servers (e.g. Render) run in UTC, so "now" must be converted before it is
// compared with the configured hours.
func LoadBusinessLocation() *time.Location {
	name := strings.TrimSpace(os.Getenv("BUSINESS_TIMEZONE"))
	if name == "" {
		name = DefaultBusinessTimezone
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Printf("invalid BUSINESS_TIMEZONE %q, falling back to %s: %v", name, DefaultBusinessTimezone, err)
		loc, _ = time.LoadLocation(DefaultBusinessTimezone)
	}
	return loc
}

// clockMinutes converts "HH:MM" to minutes after midnight.
func clockMinutes(value string) (int, bool) {
	if !clockPattern.MatchString(value) {
		return 0, false
	}
	return int(value[0]-'0')*600 + int(value[1]-'0')*60 + int(value[3]-'0')*10 + int(value[4]-'0'), true
}

// hoursFor returns the entry for day, if any.
func hoursFor(hours []models.OpeningHours, day time.Weekday) (models.OpeningHours, bool) {
	for _, h := range hours {
		if h.Day == weekdayNames[day] {
			return h, true
		}
	}
	return models.OpeningHours{}, false
}

// window returns the open and close minutes of an open day; overnight reports
// whether trading runs past midnight into the next day.
func window(h models.OpeningHours) (open, close int, overnight, ok bool) {
	if !h.IsOpen {
		return 0, 0, false, false
	}
	open, okOpen := clockMinutes(h.OpenTime)
	close, okClose := clockMinutes(h.CloseTime)
	if !okOpen || !okClose {
		return 0, 0, false, false
	}
	return open, close, close <= open, true
}

// isOpenAt reports whether t (already in the business timezone) falls inside
// the opening hours. The close time is exclusive, and a window that runs past
// midnight is honoured on the following morning.
func isOpenAt(hours []models.OpeningHours, t time.Time) bool {
	minute := t.Hour()*60 + t.Minute()

	if today, found := hoursFor(hours, t.Weekday()); found {
		if open, close, overnight, ok := window(today); ok {
			if overnight && minute >= open {
				return true
			}
			if !overnight && minute >= open && minute < close {
				return true
			}
		}
	}

	if yesterday, found := hoursFor(hours, (t.Weekday()+6)%7); found {
		if _, close, overnight, ok := window(yesterday); ok && overnight && minute < close {
			return true
		}
	}

	return false
}

// normalizeOpeningHours validates a full week of opening hours and returns it
// in Monday-first order with canonical day names. Every day must appear
// exactly once; closed days may leave their times empty.
func normalizeOpeningHours(hours []models.OpeningHours) ([]models.OpeningHours, error) {
	if len(hours) != len(weekdayNames) {
		return nil, fmt.Errorf("opening hours must list all 7 days, got %d", len(hours))
	}

	byDay := make(map[time.Weekday]models.OpeningHours, len(hours))
	for _, h := range hours {
		day, ok := parseWeekday(h.Day)
		if !ok {
			return nil, fmt.Errorf("unknown day %q in opening hours", h.Day)
		}
		if _, dup := byDay[day]; dup {
			return nil, fmt.Errorf("%s appears more than once in opening hours", weekdayNames[day])
		}

		h.Day = weekdayNames[day]
		h.OpenTime = strings.TrimSpace(h.OpenTime)
		h.CloseTime = strings.TrimSpace(h.CloseTime)

		for label, value := range map[string]string{"opening": h.OpenTime, "closing": h.CloseTime} {
			if value == "" && !h.IsOpen {
				continue
			}
			if _, ok := clockMinutes(value); !ok {
				return nil, fmt.Errorf("%s %s time %q must be in 24-hour HH:MM format", h.Day, label, value)
			}
		}
		if h.IsOpen && h.OpenTime == h.CloseTime {
			return nil, fmt.Errorf("%s opening and closing times cannot be the same", h.Day)
		}

		byDay[day] = h
	}

	normalized := make([]models.OpeningHours, 0, len(weekdayNames))
	for i := 1; i <= len(weekdayNames); i++ {
		normalized = append(normalized, byDay[time.Weekday(i%7)])
	}
	return normalized, nil
}

func parseWeekday(value string) (time.Weekday, bool) {
	value = strings.TrimSpace(value)
	for i, name := range weekdayNames {
		if strings.EqualFold(value, name) {
			return time.Weekday(i), true
		}
	}
	return 0, false
}

// normalizeSettingsUpdate trims text fields in place and rejects values the
// storefront could not display or that would not fit their columns.
func normalizeSettingsUpdate(u *models.BusinessSettingsUpdate) error {
	text := []struct {
		label    string
		value    *string
		required bool
		max      int
	}{
		{"business name", u.BusinessName, true, 200},
		{"phone number", u.PhoneNumber, true, 20},
		{"email", u.Email, true, 200},
		{"street address", u.Address, true, 500},
		{"city", u.City, true, 100},
		{"state", u.State, true, 100},
		{"postal code", u.PostalCode, false, 20},
		{"country", u.Country, true, 100},
	}
	for _, field := range text {
		if field.value == nil {
			continue
		}
		*field.value = strings.TrimSpace(*field.value)
		if field.required && *field.value == "" {
			return fmt.Errorf("%s is required", field.label)
		}
		if len([]rune(*field.value)) > field.max {
			return fmt.Errorf("%s must be at most %d characters", field.label, field.max)
		}
	}

	if u.PhoneNumber != nil && !phonePattern.MatchString(*u.PhoneNumber) {
		return fmt.Errorf("phone number must contain only digits, spaces, (), - and an optional leading +")
	}
	if u.Email != nil {
		if parsed, err := mail.ParseAddress(*u.Email); err != nil || parsed.Address != *u.Email {
			return fmt.Errorf("email must be a valid email address")
		}
	}

	if u.OpeningHours != nil {
		normalized, err := normalizeOpeningHours(*u.OpeningHours)
		if err != nil {
			return err
		}
		u.OpeningHours = &normalized
	}

	return nil
}
