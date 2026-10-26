package utils

import "testing"

func TestCleanPhone(t *testing.T) {
	cases := map[string]string{
		" +234 (801) 234-5678 ": "+2348012345678",
		"0801 234 5678":         "08012345678",
		"080+1":                 "0801",
	}
	for in, want := range cases {
		if got := CleanPhone(in); got != want {
			t.Errorf("CleanPhone(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsValidPhone(t *testing.T) {
	for value, want := range map[string]bool{
		"08012345678":    true,
		"+2348012345678": true,
		"12345":          false,
		"":               false,
	} {
		if got := IsValidPhone(value); got != want {
			t.Errorf("IsValidPhone(%q) = %t, want %t", value, got, want)
		}
	}
}

func TestSamePhone(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"+234 801 234 5678", "08012345678", true},
		{"2348012345678", "0801-234-5678", true},
		{"08012345678", "08012345679", false},
		{"", "", false},
	}
	for _, tc := range cases {
		if got := SamePhone(tc.a, tc.b); got != tc.want {
			t.Errorf("SamePhone(%q, %q) = %t, want %t", tc.a, tc.b, got, tc.want)
		}
	}
}
