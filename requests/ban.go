package requests

// Ban struct
type Ban struct {
	ChatId         int64 `json:"chat_id"`
	UserId         int64 `json:"user_id"`
	UntilDate      int64 `json:"until_date,omitempty"`
	RevokeMessages bool  `json:"revoke_messages,omitempty"`
}

// UnBan struct
type UnBan struct {
	ChatId       int64 `json:"chat_id"`
	UserId       int64 `json:"user_id"`
	OnlyIfBanned bool  `json:"only_if_banned,omitempty"`
}
