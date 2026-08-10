package requests

import (
	"github.com/bytedance/sonic"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type EditMessage struct {
	ChatId             int64               `json:"chat_id"`
	MessageId          int64               `json:"message_id"`
	Text               string              `json:"text,omitempty"`
	Caption            string              `json:"caption,omitempty"`
	ParseMode          string              `json:"parse_mode,omitempty"`
	Media              *InputMedia         `json:"media,omitempty"`
	ReplyMarkup        *models.ReplyMarkup `json:"reply_markup,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Buttons            [][]map[string]any  `json:"-"`
	ButtonType         string              `json:"-"`
}

func (p *EditMessage) MarshalJSON() ([]byte, error) {
	data := *p

	if data.ParseMode == "" {
		data.ParseMode = "html"
	}

	if len(data.Buttons) >= 1 {
		if data.ButtonType == "keyboard" {
			data.ReplyMarkup = utils.SetKeyboard(data.Buttons, true)
		} else {
			data.ReplyMarkup = utils.SetInlineKeyboard(data.Buttons)
		}
	}

	return sonic.Marshal(data)
}

type InputMedia struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}
