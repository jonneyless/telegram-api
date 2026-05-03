package models

// ReplyMarkup struct
type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton struct
type InlineKeyboardButton struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}
