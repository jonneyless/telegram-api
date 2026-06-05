package requests

// Tag struct
type Tag struct {
	ChatId int64  `json:"chat_id"`
	UserId int64  `json:"user_id"`
	Tag    string `json:"tag"`
}
