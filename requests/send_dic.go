package requests

import (
	"github.com/bytedance/sonic"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type SendDice struct {
	ChatId              int64               `json:"chat_id"`
	Emoji               string              `json:"emoji,omitempty"`
	DisableNotification bool                `json:"disable_notification,omitempty"`
	ProtectContent      bool                `json:"protect_content,omitempty"`
	ReplyParameters     *ReplyParameters    `json:"reply_parameters,omitempty"`
	ReplyMarkup         *models.ReplyMarkup `json:"reply_markup,omitempty"`
	Buttons             [][]map[string]any  `json:"-"`
	ButtonType          string              `json:"-"`
}

func (p *SendDice) MarshalJSON() ([]byte, error) {
	data := *p

	if data.Emoji == "" {
		data.Emoji = "🎲"
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
