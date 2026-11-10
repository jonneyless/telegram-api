package requests

// PinMessage Message struct
type PinMessage struct {
	ChatId              int64 `json:"chat_id"`
	MessageId           int64 `json:"message_id,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
}
