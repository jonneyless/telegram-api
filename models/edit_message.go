package models

// EditedMessage struct
type EditedMessage struct {
	MessageID   int64       `json:"message_id"`
	From        From        `json:"from"`
	Chat        Chat        `json:"chat"`
	Date        int64       `json:"date"`
	EditDate    int64       `json:"edit_date"`
	Text        string      `json:"text"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
	Entities    []Entities  `json:"entities"`
}
