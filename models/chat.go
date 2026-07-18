package models

// Chat struct
type Chat struct {
	ID               int64  `json:"id"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	Username         string `json:"username"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	IsForum          bool   `json:"is_forum"`
	IsDirectMessages bool   `json:"is_direct_messages"`
}

// IsPrivate 私聊
func (m *Chat) IsPrivate() bool {
	return m.Type == "private"
}

type ChatFull struct {
	ID                          int64    `json:"id"`
	Type                        string   `json:"type"`
	Title                       string   `json:"title"`
	Username                    string   `json:"username"`
	FirstName                   string   `json:"first_name"`
	LastName                    string   `json:"last_name"`
	Bio                         string   `json:"bio"`
	Description                 string   `json:"description"`
	InviteLink                  string   `json:"invite_link"`
	HasVisibleHistory           bool     `json:"has_visible_history"`
	AllMembersAreAdministrators bool     `json:"all_members_are_administrators"`
	ActiveUsernames             []string `json:"active_usernames,omitempty"`
	Photo                       *Avatar  `json:"photo,omitempty"`
}
