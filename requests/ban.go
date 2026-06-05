package requests

// Ban struct
type Ban struct {
	ChatId         int64 `json:"chat_id"`
	UserId         int64 `json:"user_id"`
	UntilDate      int64 `json:"until_date"`
	RevokeMessages bool  `json:"revoke_messages"`
}

// UnBan struct
type UnBan struct {
	ChatId       int64 `json:"chat_id"`
	UserId       int64 `json:"user_id"`
	OnlyIfBanned bool  `json:"only_if_banned"`
}
