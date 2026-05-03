package requests

type CreateChatInviteLink struct {
	ChatId             int64  `json:"chat_id"`
	Name               string `json:"name"`
	ExpireDate         int64  `json:"expire_date"`
	MemberLimit        int64  `json:"member_limit"`
	CreatesJoinRequest bool   `json:"creates_join_request"`
}
