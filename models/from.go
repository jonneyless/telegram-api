package models

import "strings"

// From struct
type From struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	IsPremium    bool   `json:"is_premium"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

func (form *From) FullName() string {
	return strings.TrimSpace(form.FirstName + " " + form.LastName)
}
