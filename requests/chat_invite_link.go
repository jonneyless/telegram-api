package requests

type CreateChatInviteLink struct {
	ChatId             int64  `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int64  `json:"expire_date,omitempty"`
	MemberLimit        int64  `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

type RevokeChatInviteLink struct {
	ChatId     int64  `json:"chat_id"`
	InviteLink string `json:"invite_link"`
}
