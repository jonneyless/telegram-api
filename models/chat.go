package models

// Chat struct
type Chat struct {
	ID                          int64  `json:"id"`
	FirstName                   string `json:"first_name"`
	LastName                    string `json:"last_name"`
	Username                    string `json:"username"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	InviteLink                  string `json:"invite_link"`
	HasVisibleHistory           bool   `json:"has_visible_history"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}
