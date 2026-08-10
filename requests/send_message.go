package requests

import (
	"github.com/bytedance/sonic"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type SendMessage struct {
	ChatId              int64               `json:"chat_id"`
	Text                string              `json:"text,omitempty"`
	Photo               string              `json:"photo,omitempty"`
	Video               string              `json:"video,omitempty"`
	Audio               string              `json:"audio,omitempty"`
	Document            string              `json:"document,omitempty"`
	ParseMode           string              `json:"parse_mode,omitempty"`
	ReplyMarkup         *models.ReplyMarkup `json:"reply_markup,omitempty"`
	DisableNotification bool                `json:"disable_notification,omitempty"`
	LinkPreviewOptions  *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	ReplyParameters     *ReplyParameters    `json:"reply_parameters,omitempty"`
	Buttons             [][]map[string]any  `json:"-"`
	ButtonType          string              `json:"-"`
}

func (p *SendMessage) MarshalJSON() ([]byte, error) {
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

type LinkPreviewOptions struct {
	IsDisabled bool   `json:"is_disabled"`
	Url        string `json:"url,omitempty"`
}

type ReplyParameters struct {
	MessageId int64 `json:"message_id"`
	ChatId    int64 `json:"chat_id,omitempty"`
}
