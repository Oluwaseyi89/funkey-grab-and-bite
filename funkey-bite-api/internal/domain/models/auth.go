package models

type AuthContext struct {
	UserID  int    `json:"user_id"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"is_admin"`
}
