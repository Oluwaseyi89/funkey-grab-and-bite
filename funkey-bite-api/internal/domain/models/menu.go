package models

import (
	"encoding/json"
	"time"
)

type MenuCategory struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	DisplayOrder int       `json:"displayOrder" db:"display_order"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type MenuItem struct {
	ID              int              `json:"id" db:"id"`
	CategoryID      int              `json:"categoryId" db:"category_id"`
	Name            string           `json:"name" db:"name"`
	Description     string           `json:"description" db:"description"`
	Price           float64          `json:"price" db:"price"`
	ImageURL        string           `json:"imageUrl" db:"image_url"`
	IsAvailable     bool             `json:"isAvailable" db:"is_available"`
	IsPreOrder      bool             `json:"isPreOrder" db:"is_pre_order"`
	PreparationTime int              `json:"preparationTime" db:"preparation_time"`
	Tags            []string         `json:"tags" db:"tags"`
	NutritionalInfo *NutritionalInfo `json:"nutritionalInfo,omitempty"`
	CreatedAt       time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time        `json:"updatedAt,omitempty" db:"updated_at"`
}

// imageURLResolver turns a stored image value (an object key) into the URL
// clients load. It is the identity until storage is configured at startup.
var imageURLResolver = func(value string) string { return value }

// SetImageURLResolver installs the key-to-URL mapping used when menu items are
// encoded, so every endpoint returning a MenuItem serves full image URLs.
func SetImageURLResolver(resolve func(string) string) {
	if resolve != nil {
		imageURLResolver = resolve
	}
}

func (m MenuItem) MarshalJSON() ([]byte, error) {
	type menuItemJSON MenuItem // drops this method, avoiding recursion
	out := menuItemJSON(m)
	out.ImageURL = imageURLResolver(m.ImageURL)
	return json.Marshal(out)
}

type NutritionalInfo struct {
	Calories int `json:"calories"`
	Protein  int `json:"protein"`
	Carbs    int `json:"carbs"`
	Fat      int `json:"fat"`
}
