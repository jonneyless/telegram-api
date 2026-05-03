package models

// ChannelPost struct
type ChannelPost struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}
