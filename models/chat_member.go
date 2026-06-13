package models

// ChatMember struct
type ChatMember struct {
	MessageID     int64       `json:"message_id"`
	From          From        `json:"from"`
	Chat          Chat        `json:"chat"`
	Date          int64       `json:"date"`
	OldChatMember Member      `json:"old_chat_member"`
	NewChatMember Member      `json:"new_chat_member"`
	InviteLink    *InviteLink `json:"invite_link"`
}

// IsLeave 离群
func (m *ChatMember) IsLeave() bool {
	return m.NewChatMember.IsLeft()
}

// IsNewMember 入群
func (m *ChatMember) IsNewMember() bool {
	if !m.OldChatMember.IsLeft() && !m.OldChatMember.IsBanned() {
		return false
	}
	return m.NewChatMember.IsChatMember() || m.NewChatMember.IsAdministrator()
}

// IsPromoteToAdministrator 设为管理
func (m *ChatMember) IsPromoteToAdministrator() bool {
	return m.NewChatMember.IsAdministrator()
}

// IsRevokeAdministrator 取消管理
func (m *ChatMember) IsRevokeAdministrator() bool {
	return m.OldChatMember.IsAdministrator() && m.NewChatMember.IsChatMember()
}

// IsRestricted 被禁言
func (m *ChatMember) IsRestricted() bool {
	return m.NewChatMember.IsRestricted()
}
