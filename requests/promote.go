package requests

type BasePromote struct {
	ChatId      int64 `json:"chat_id"`
	UserId      int64 `json:"user_id"`
	IsAnonymous bool  `json:"is_anonymous"`
}

type PromoteChat struct {
	BasePromote
}

type PromoteSuperGroup struct {
	BasePromote
	CanManageChat           bool `json:"can_manage_chat"`
	CanDeleteMessages       bool `json:"can_delete_messages"`
	CanManageVideoChats     bool `json:"can_manage_video_chats"`
	CanRestrictMembers      bool `json:"can_restrict_members"`
	CanPromoteMembers       bool `json:"can_promote_members"`
	CanChangeInfo           bool `json:"can_change_info"`
	CanInviteUsers          bool `json:"can_invite_users"`
	CanPostStories          bool `json:"can_post_stories"`
	CanEditStories          bool `json:"can_edit_stories"`
	CanDeleteStories        bool `json:"can_delete_stories"`
	CanPinMessages          bool `json:"can_pin_messages"`
	CanManageTopics         bool `json:"can_manage_topics"`
	CanManageDirectMessages bool `json:"can_manage_direct_messages"`
	CanManageTags           bool `json:"can_manage_tags"`
}

type PromoteChannel struct {
	BasePromote
	CanChangeInfo    bool `json:"can_change_info"`
	CanInviteUsers   bool `json:"can_invite_users"`
	CanPostStories   bool `json:"can_post_stories"`
	CanEditStories   bool `json:"can_edit_stories"`
	CanDeleteStories bool `json:"can_delete_stories"`
	CanPostMessages  bool `json:"can_post_messages"`
	CanEditMessages  bool `json:"can_edit_messages"`
	CanPinMessages   bool `json:"can_pin_messages"`
}

type Promotable interface {
	GetChatId() int64
	GetUserId() int64
}

func (b BasePromote) GetChatId() int64 { return b.ChatId }
func (b BasePromote) GetUserId() int64 { return b.UserId }
