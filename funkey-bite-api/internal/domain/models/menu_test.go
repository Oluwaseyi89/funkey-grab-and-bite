package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMenuItemJSONResolvesImageURL(t *testing.T) {
	original := imageURLResolver
	t.Cleanup(func() { imageURLResolver = original })

	SetImageURLResolver(func(value string) string { return "https://cdn.example.com/" + value })

	encoded, err := json.Marshal([]MenuItem{{ID: 1, Name: "Shawarma", ImageURL: "menu-items/a.webp"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), `"imageUrl":"https://cdn.example.com/menu-items/a.webp"`) {
		t.Fatalf("image key was not resolved: %s", encoded)
	}
	if !strings.Contains(string(encoded), `"name":"Shawarma"`) {
		t.Fatalf("other fields must be encoded unchanged: %s", encoded)
	}
}
