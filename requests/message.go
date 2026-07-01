package requests

// Message struct
type Message struct {
	ChatId     int64    `json:"chat_id"`
	MessageId  *int64   `json:"message_id,omitempty"`
	MessageIds *[]int64 `json:"message_ids,omitempty"`
}
