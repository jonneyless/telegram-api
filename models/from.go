package models

import (
	"fmt"
	"strings"
)

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

func (f *From) FullName() string {
	return strings.TrimSpace(f.FirstName + " " + f.LastName)
}

func (f *From) GetMention() string {
	if f.Username != "" {
		return "@" + f.Username
	}

	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, f.ID, f.FullName())
}
