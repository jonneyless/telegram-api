package utils

import (
	"github.com/jonneyless/telegram-api/models"
	"github.com/spf13/cast"
)

func SetInlineKeyboard(data [][]map[string]any) *models.ReplyMarkup {
	inlineKeyboard := make([][]models.InlineKeyboardButton, 0)

	for _, row := range data {
		inlineKeyboardRow := make([]models.InlineKeyboardButton, 0)
		for _, buttonMap := range row {
			if buttonMap == nil {
				continue
			}

			var button models.InlineKeyboardButton

			if val, ok := buttonMap["text"]; ok {
				if s, ok := val.(string); ok {
					button.Text = s
				} else {
					continue
				}
			} else {
				continue
			}

			if val, ok := buttonMap["url"]; ok {
				if s, ok := val.(string); ok {
					button.URL = s
				} else {
					continue
				}
			}

			if val, ok := buttonMap["callback_data"]; ok {
				if s, ok := val.(string); ok {
					button.CallbackData = s
				} else {
					continue
				}
			}

			inlineKeyboardRow = append(inlineKeyboardRow, button)
		}

		inlineKeyboard = append(inlineKeyboard, inlineKeyboardRow)
	}

	if inlineKeyboard != nil {
		return &models.ReplyMarkup{
			InlineKeyboard: inlineKeyboard,
		}
	}

	return nil
}

func SetKeyboard(data [][]map[string]any, resize bool) *models.ReplyMarkup {
	keyboard := make([][]models.KeyboardButton, 0)

	for _, row := range data {
		keyboardRow := make([]models.KeyboardButton, 0)
		for _, buttonMap := range row {
			if buttonMap == nil {
				continue
			}

			var button models.KeyboardButton

			if val, ok := buttonMap["text"]; ok {
				if s, ok := val.(string); ok {
					button.Text = s
				} else {
					continue
				}
			} else {
				continue
			}

			if val, ok := buttonMap["request_chat"]; ok {
				if s, ok := val.(map[string]any); ok {
					button.RequestChat = convertKeyboardButtonRequestChat(s)
				} else {
					continue
				}
			}

			keyboardRow = append(keyboardRow, button)
		}

		keyboard = append(keyboard, keyboardRow)
	}

	if keyboard != nil {
		return &models.ReplyMarkup{
			Keyboard:       keyboard,
			ResizeKeyboard: resize,
		}
	}

	return nil
}

func convertKeyboardButtonRequestChat(m map[string]any) *models.KeyboardButtonRequestChat {
	requestChat := &models.KeyboardButtonRequestChat{}

	if val, ok := m["request_id"]; ok {
		requestChat.RequestId = cast.ToInt32(val)
	}

	if val, ok := m["chat_is_channel"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.ChatIsChannel = s
		}
	}

	if val, ok := m["chat_is_forum"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.ChatIsForum = s
		}
	}

	if val, ok := m["chat_has_username"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.ChatHasUsername = s
		}
	}

	if val, ok := m["chat_is_created"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.ChatIsCreated = s
		}
	}

	if val, ok := m["bot_is_member"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.BotIsMember = s
		}
	}

	if val, ok := m["request_title"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.RequestTitle = s
		}
	}

	if val, ok := m["request_username"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.RequestUsername = s
		}
	}

	if val, ok := m["request_photo"]; ok {
		if s, ok := val.(bool); ok {
			requestChat.RequestPhoto = s
		}
	}

	return requestChat
}
