package models

type Member struct {
	User        From    `json:"user"`
	Status      string  `json:"status"`
	CustomTitle *string `json:"custom_title,omitempty"`
	Tag         *string `json:"tag,omitempty"`
	UntilDate   *int64  `json:"until_date,omitempty"`
	// status == administrator
	CanBeEdited             *bool `json:"can_be_edited,omitempty"`
	CanManageChat           *bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages       *bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats     *bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers      *bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers       *bool `json:"can_promote_members,omitempty"`
	CanPostStories          *bool `json:"can_post_stories,omitempty"`
	CanEditStories          *bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories        *bool `json:"can_delete_stories,omitempty"`
	CanPostMessages         *bool `json:"can_post_messages,omitempty"`
	CanEditMessages         *bool `json:"can_edit_messages,omitempty"`
	CanManageDirectMessages *bool `json:"can_manage_direct_messages,omitempty"`
	CanManageTags           *bool `json:"can_manage_tags,omitempty"`
	IsAnonymous             *bool `json:"is_anonymous,omitempty"`
	// status == restricted
	IsMember              *bool `json:"is_member,omitempty"`
	CanSendMessages       *bool `json:"can_send_messages,omitempty"`
	CanSendAudios         *bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      *bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         *bool `json:"can_send_photos,omitempty"`
	CanSendVideos         *bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     *bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     *bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          *bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  *bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews *bool `json:"can_add_web_page_previews,omitempty"`
	CanReactToMessages    *bool `json:"can_react_to_messages,omitempty"`
	CanEditTag            *bool `json:"can_edit_tag,omitempty"`
	// status == restricted or status == administrator
	CanChangeInfo   *bool `json:"can_change_info,omitempty"`
	CanInviteUsers  *bool `json:"can_invite_users,omitempty"`
	CanPinMessages  *bool `json:"can_pin_messages,omitempty"`
	CanManageTopics *bool `json:"can_manage_topics,omitempty"`
}

func (m *Member) IsOwner() bool {
	return m.Status == "creator"
}

func (m *Member) IsAdministrator() bool {
	return m.Status == "administrator"
}

func (m *Member) IsChatMember() bool {
	return m.Status == "member"
}

func (m *Member) IsRestricted() bool {
	return m.Status == "restricted"
}

func (m *Member) IsBanned() bool {
	return m.Status == "kicked"
}

func (m *Member) IsLeft() bool {
	return m.Status == "left" || m.Status == ""
}

func (m *Member) Permissions() map[string]bool {
	return map[string]bool{
		"can_be_edited":              m.GetCanBeEdited(),
		"can_manage_chat":            m.GetCanManageChat(),
		"can_delete_messages":        m.GetCanDeleteMessages(),
		"can_manage_video_chats":     m.GetCanManageVideoChats(),
		"can_restrict_members":       m.GetCanRestrictMembers(),
		"can_promote_members":        m.GetCanPromoteMembers(),
		"can_change_info":            m.GetCanChangeInfo(),
		"can_invite_users":           m.GetCanInviteUsers(),
		"can_post_stories":           m.GetCanPostStories(),
		"can_edit_stories":           m.GetCanEditStories(),
		"can_delete_stories":         m.GetCanDeleteStories(),
		"can_post_messages":          m.GetCanPostMessages(),
		"can_edit_messages":          m.GetCanEditMessages(),
		"can_pin_messages":           m.GetCanPinMessages(),
		"can_manage_topics":          m.GetCanManageTopics(),
		"can_manage_direct_messages": m.GetCanManageDirectMessages(),
		"can_manage_tags":            m.GetCanManageTags(),
		"is_anonymous":               m.GetIsAnonymous(),
	}
}

func (m *Member) Restricted() map[string]bool {
	return map[string]bool{
		"can_send_messages":         m.GetCanSendMessages(),
		"can_send_audios":           m.GetCanSendAudios(),
		"can_send_documents":        m.GetCanSendDocuments(),
		"can_send_photos":           m.GetCanSendPhotos(),
		"can_send_videos":           m.GetCanSendVideos(),
		"can_send_video_notes":      m.GetCanSendVideoNotes(),
		"can_send_voice_notes":      m.GetCanSendVoiceNotes(),
		"can_send_polls":            m.GetCanSendPolls(),
		"can_send_other_messages":   m.GetCanSendOtherMessages(),
		"can_add_web_page_previews": m.GetCanAddWebPagePreviews(),
		"can_react_to_messages":     m.GetCanReactToMessages(),
		"can_edit_tag":              m.GetCanEditTag(),
		"can_change_info":           m.GetCanChangeInfo(),
		"can_invite_users":          m.GetCanInviteUsers(),
		"can_pin_messages":          m.GetCanPinMessages(),
		"can_manage_topics":         m.GetCanManageTopics(),
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

func (m *Member) GetCanDeleteMessages() bool {
	if m != nil && m.CanDeleteMessages != nil {
		return *m.CanDeleteMessages
	}

	return false
}

func (m *Member) GetCanManageVideoChats() bool {
	if m != nil && m.CanManageVideoChats != nil {
		return *m.CanManageVideoChats
	}

	return false
}

func (m *Member) GetCanRestrictMembers() bool {
	if m != nil && m.CanRestrictMembers != nil {
		return *m.CanRestrictMembers
	}

	return false
}

func (m *Member) GetCanPromoteMembers() bool {
	if m != nil && m.CanPromoteMembers != nil {
		return *m.CanPromoteMembers
	}

	return false
}

func (m *Member) GetCanChangeInfo() bool {
	if m != nil && m.CanChangeInfo != nil {
		return *m.CanChangeInfo
	}

	return false
}

func (m *Member) GetCanInviteUsers() bool {
	if m != nil && m.CanInviteUsers != nil {
		return *m.CanInviteUsers
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

func (m *Member) GetCanPostMessages() bool {
	if m != nil && m.CanPostMessages != nil {
		return *m.CanPostMessages
	}

	return false
}

func (m *Member) GetCanEditMessages() bool {
	if m != nil && m.CanEditMessages != nil {
		return *m.CanEditMessages
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

func (m *Member) GetCanManageDirectMessages() bool {
	if m != nil && m.CanManageDirectMessages != nil {
		return *m.CanManageDirectMessages
	}

	return false
}

func (m *Member) GetCanManageTags() bool {
	if m != nil && m.CanManageTags != nil {
		return *m.CanManageTags
	}

	return false
}

func (m *Member) GetIsAnonymous() bool {
	if m != nil && m.IsAnonymous != nil {
		return *m.IsAnonymous
	}

	return false
}

func (m *Member) GetCanSendMessages() bool {
	if m != nil && m.CanSendMessages != nil {
		return *m.CanSendMessages
	}
	return false
}

func (m *Member) GetCanSendAudios() bool {
	if m != nil && m.CanSendAudios != nil {
		return *m.CanSendAudios
	}
	return false
}

func (m *Member) GetCanSendDocuments() bool {
	if m != nil && m.CanSendDocuments != nil {
		return *m.CanSendDocuments
	}
	return false
}

func (m *Member) GetCanSendPhotos() bool {
	if m != nil && m.CanSendPhotos != nil {
		return *m.CanSendPhotos
	}
	return false
}

func (m *Member) GetCanSendVideos() bool {
	if m != nil && m.CanSendVideos != nil {
		return *m.CanSendVideos
	}
	return false
}

func (m *Member) GetCanSendVideoNotes() bool {
	if m != nil && m.CanSendVideoNotes != nil {
		return *m.CanSendVideoNotes
	}
	return false
}

func (m *Member) GetCanSendVoiceNotes() bool {
	if m != nil && m.CanSendVoiceNotes != nil {
		return *m.CanSendVoiceNotes
	}
	return false
}

func (m *Member) GetCanSendPolls() bool {
	if m != nil && m.CanSendPolls != nil {
		return *m.CanSendPolls
	}
	return false
}

func (m *Member) GetCanSendOtherMessages() bool {
	if m != nil && m.CanSendOtherMessages != nil {
		return *m.CanSendOtherMessages
	}
	return false
}

func (m *Member) GetCanAddWebPagePreviews() bool {
	if m != nil && m.CanAddWebPagePreviews != nil {
		return *m.CanAddWebPagePreviews
	}
	return false
}

func (m *Member) GetCanReactToMessages() bool {
	if m != nil && m.CanReactToMessages != nil {
		return *m.CanReactToMessages
	}
	return false
}

func (m *Member) GetCanEditTag() bool {
	if m != nil && m.CanEditTag != nil {
		return *m.CanEditTag
	}
	return false
}
