package models

import (
	"encoding/json"
	"time"
)

// TeamMember is a person shown in the team section of the contact page.
type TeamMember struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Bio          string    `json:"bio"`
	PhotoURL     *string   `json:"photoUrl,omitempty"`
	DisplayOrder int       `json:"displayOrder"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// MarshalJSON serves uploaded photos as full URLs, like menu item images.
func (m TeamMember) MarshalJSON() ([]byte, error) {
	type teamMemberJSON TeamMember // drops this method, avoiding recursion
	out := teamMemberJSON(m)
	if m.PhotoURL != nil {
		resolved := imageURLResolver(*m.PhotoURL)
		out.PhotoURL = &resolved
	}
	return json.Marshal(out)
}

// TeamMemberInput creates or replaces a team member.
type TeamMemberInput struct {
	Name         string  `json:"name" validate:"required,min=2,max=120"`
	Role         string  `json:"role" validate:"required,min=2,max=120"`
	Bio          string  `json:"bio" validate:"max=500"`
	PhotoURL     *string `json:"photoUrl,omitempty" validate:"omitempty,max=500"`
	DisplayOrder int     `json:"displayOrder" validate:"min=0,max=1000"`
	IsActive     bool    `json:"isActive"`
}
