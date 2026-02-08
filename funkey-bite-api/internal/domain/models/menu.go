package models

import "time"

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

type NutritionalInfo struct {
	Calories int `json:"calories"`
	Protein  int `json:"protein"`
	Carbs    int `json:"carbs"`
	Fat      int `json:"fat"`
}
