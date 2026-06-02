package models

type Member struct {
	User                From    `json:"user"`
	Status              string  `json:"status"`
	CustomTitle         *string `json:"custom_title"`
	Tag                 *string `json:"tag"`
	CanBeEdited         *bool   `json:"can_be_edited"`
	CanManageChat       *bool   `json:"can_manage_chat"`
	CanChangeInfo       *bool   `json:"can_change_info"`
	CanDeleteMessages   *bool   `json:"can_delete_messages"`
	CanInviteUsers      *bool   `json:"can_invite_users"`
	CanRestrictMembers  *bool   `json:"can_restrict_members"`
	CanPinMessages      *bool   `json:"can_pin_messages"`
	CanManageTopics     *bool   `json:"can_manage_topics"`
	CanPromoteMembers   *bool   `json:"can_promote_members"`
	CanManageVideoChats *bool   `json:"can_manage_video_chats"`
	CanPostStories      *bool   `json:"can_post_stories"`
	CanEditStories      *bool   `json:"can_edit_stories"`
	CanDeleteStories    *bool   `json:"can_delete_stories"`
	IsAnonymous         *bool   `json:"is_anonymous"`
	CanManageVoiceChats *bool   `json:"can_manage_voice_chats"`
	IsMember            *bool   `json:"is_member"`
}

func (m *Member) Permissions() map[string]bool {
	return map[string]bool{
		"can_be_edited":          m.GetCanBeEdited(),
		"can_manage_chat":        m.GetCanManageChat(),
		"can_change_info":        m.GetCanChangeInfo(),
		"can_delete_messages":    m.GetCanDeleteMessages(),
		"can_invite_users":       m.GetCanInviteUsers(),
		"can_restrict_members":   m.GetCanRestrictMembers(),
		"can_pin_messages":       m.GetCanPinMessages(),
		"can_manage_topics":      m.GetCanManageTopics(),
		"can_promote_members":    m.GetCanPromoteMembers(),
		"can_manage_video_chats": m.GetCanManageVideoChats(),
		"can_post_stories":       m.GetCanPostStories(),
		"can_edit_stories":       m.GetCanEditStories(),
		"can_delete_stories":     m.GetCanDeleteStories(),
		"is_anonymous":           m.GetIsAnonymous(),
		"can_manage_voice_chats": m.GetCanManageVoiceChats(),
	}
}

func (m *Member) GetCanBeEdited() bool {
	if m != nil && m.CanBeEdited != nil {
		return *m.CanBeEdited
	}

	return false
}

func (m *Member) GetCanManageChat() bool {
	if m != nil && m.CanManageChat != nil {
		return *m.CanManageChat
	}

	return false
}

func (m *Member) GetCanChangeInfo() bool {
	if m != nil && m.CanChangeInfo != nil {
		return *m.CanChangeInfo
	}

	return false
}

func (m *Member) GetCanDeleteMessages() bool {
	if m != nil && m.CanDeleteMessages != nil {
		return *m.CanDeleteMessages
	}

	return false
}

func (m *Member) GetCanInviteUsers() bool {
	if m != nil && m.CanInviteUsers != nil {
		return *m.CanInviteUsers
	}

	return false
}

func (m *Member) GetCanRestrictMembers() bool {
	if m != nil && m.CanRestrictMembers != nil {
		return *m.CanRestrictMembers
	}

	return false
}

func (m *Member) GetCanPinMessages() bool {
	if m != nil && m.CanPinMessages != nil {
		return *m.CanPinMessages
	}

	return false
}

func (m *Member) GetCanManageTopics() bool {
	if m != nil && m.CanManageTopics != nil {
		return *m.CanManageTopics
	}

	return false
}

func (m *Member) GetCanPromoteMembers() bool {
	if m != nil && m.CanPromoteMembers != nil {
		return *m.CanPromoteMembers
	}

	return false
}

func (m *Member) GetCanManageVideoChats() bool {
	if m != nil && m.CanManageVideoChats != nil {
		return *m.CanManageVideoChats
	}

	return false
}

func (m *Member) GetCanPostStories() bool {
	if m != nil && m.CanPostStories != nil {
		return *m.CanPostStories
	}

	return false
}

func (m *Member) GetCanEditStories() bool {
	if m != nil && m.CanEditStories != nil {
		return *m.CanEditStories
	}

	return false
}

func (m *Member) GetCanDeleteStories() bool {
	if m != nil && m.CanDeleteStories != nil {
		return *m.CanDeleteStories
	}

	return false
}

func (m *Member) GetIsAnonymous() bool {
	if m != nil && m.IsAnonymous != nil {
		return *m.IsAnonymous
	}

	return false
}

func (m *Member) GetCanManageVoiceChats() bool {
	if m != nil && m.CanManageVoiceChats != nil {
		return *m.CanManageVoiceChats
	}

	return false
}
