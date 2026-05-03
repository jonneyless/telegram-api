package requests

// Member struct
type Member struct {
	Chat
	UserId int64 `json:"user_id"`
}
