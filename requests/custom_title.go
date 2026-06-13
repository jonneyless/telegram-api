package requests

// CustomTitle struct
type CustomTitle struct {
	ChatId      int64  `json:"chat_id"`
	UserId      int64  `json:"user_id"`
	CustomTitle string `json:"custom_title"`
}
