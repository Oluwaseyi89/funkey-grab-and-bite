package repository

import "testing"

func TestClockString(t *testing.T) {
	for in, want := range map[string]string{
		"0000-01-01T14:30:00Z": "14:30",
		"09:05:00":             "09:05",
	} {
		if got := clockString(in); got != want {
			t.Errorf("clockString(%q) = %q, want %q", in, got, want)
		}
	}
}
