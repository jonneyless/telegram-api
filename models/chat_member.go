package models

// ChatMember struct
type ChatMember struct {
	From          From        `json:"from"`
	Chat          Chat        `json:"chat"`
	Date          int64       `json:"date"`
	OldChatMember Member      `json:"old_chat_member"`
	NewChatMember Member      `json:"new_chat_member"`
	InviteLink    *InviteLink `json:"invite_link"`
}
