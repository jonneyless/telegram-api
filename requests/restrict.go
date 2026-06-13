package requests

// Restrict struct
type Restrict struct {
	ChatId      int64           `json:"chat_id"`
	UserId      int64           `json:"user_id"`
	Permissions ChatPermissions `json:"permissions"`
	UntilDate   int64           `json:"until_date"`
}
