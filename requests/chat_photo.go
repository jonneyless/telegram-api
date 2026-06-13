package requests

type ChatPhoto struct {
	ChatId any    `json:"chat_id"`
	Photo  []byte `json:"photo"`
}
