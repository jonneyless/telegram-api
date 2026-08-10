package models

// ReplyMarkup struct
type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
	Keyboard       [][]KeyboardButton       `json:"keyboard,omitempty"`
	IsPersistent   bool                     `json:"is_persistent,omitempty"`
	ResizeKeyboard bool                     `json:"resize_keyboard,omitempty"`
}

// KeyboardButton struct
type KeyboardButton struct {
	Text        string                     `json:"text"`
	RequestChat *KeyboardButtonRequestChat `json:"request_chat,omitempty"`
}

type KeyboardButtonRequestChat struct {
	RequestId       int32 `json:"request_id"`
	ChatIsChannel   bool  `json:"chat_is_channel"`
	ChatIsForum     bool  `json:"chat_is_forum,omitempty"`
	ChatHasUsername bool  `json:"chat_has_username,omitempty"`
	ChatIsCreated   bool  `json:"chat_is_created,omitempty"`
	BotIsMember     bool  `json:"bot_is_member,omitempty"`
	RequestTitle    bool  `json:"request_title,omitempty"`
	RequestUsername bool  `json:"request_username,omitempty"`
	RequestPhoto    bool  `json:"request_photo,omitempty"`
}

// InlineKeyboardButton struct
type InlineKeyboardButton struct {
	Text                        string                       `json:"text"`
	URL                         string                       `json:"url,omitempty"`
	CallbackData                string                       `json:"callback_data,omitempty"`
	SwitchInlineQueryChosenChat *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
}

type SwitchInlineQueryChosenChat struct {
	Query             string `json:"query,omitempty"`
	AllowUserChats    bool   `json:"allow_user_chats,omitempty"`
	AllowBotChats     bool   `json:"allow_bot_chats,omitempty"`
	AllowGroupChats   bool   `json:"allow_group_chats,omitempty"`
	AllowChannelChats bool   `json:"allow_channel_chats,omitempty"`
}
