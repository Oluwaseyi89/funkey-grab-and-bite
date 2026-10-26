package utils

import (
	"regexp"
	"strings"
)

var cleanPhonePattern = regexp.MustCompile(`^\+?[0-9]{7,15}$`)

// CleanPhone strips spaces and punctuation from a phone number, keeping a
// leading "+", so the same number is stored the same way however it was typed.
func CleanPhone(raw string) string {
	raw = strings.TrimSpace(raw)
	var b strings.Builder
	for i, r := range raw {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '+' && i == 0:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// IsValidPhone reports whether a cleaned phone number is plausible.
func IsValidPhone(cleaned string) bool {
	return cleanPhonePattern.MatchString(cleaned)
}

// SamePhone reports whether two phone numbers are the same, ignoring
// formatting and the Nigerian "+234" versus leading "0" forms.
func SamePhone(a, b string) bool {
	ca, cb := canonicalPhone(a), canonicalPhone(b)
	return ca != "" && ca == cb
}

func canonicalPhone(raw string) string {
	digits := strings.TrimPrefix(CleanPhone(raw), "+")
	if strings.HasPrefix(digits, "234") && len(digits) == 13 {
		return "0" + digits[3:]
	}
	return digits
}
