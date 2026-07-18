package requests

import (
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type SendDice struct {
	ChatId              int64               `json:"chat_id"`
	Emoji               *string             `json:"emoji,omitempty"`
	DisableNotification *bool               `json:"disable_notification,omitempty"`
	ProtectContent      *bool               `json:"protect_content,omitempty"`
	ReplyParameters     *ReplyParameters    `json:"reply_parameters,omitempty"`
	ReplyMarkup         *models.ReplyMarkup `json:"reply_markup,omitempty"`
}

func (p *SendDice) GetParams() map[string]any {
	params := map[string]any{
		"chat_id":              p.ChatId,
		"emoji":                "🎲",
		"disable_notification": p.DisableNotification,
		"protect_content":      p.ProtectContent,
	}

	if p.Emoji != nil {
		params["emoji"] = *p.Emoji
	}

	if p.ReplyMarkup != nil {
		params["reply_markup"] = utils.Struct2Map(p.ReplyMarkup)
	}

	if p.ReplyParameters != nil {
		replyParameters := map[string]any{
			"message_id": p.ReplyParameters.MessageId,
		}
		if p.ReplyParameters.ChatId != nil {
			replyParameters["chat_id"] = p.ReplyParameters.ChatId
		}

		params["reply_parameters"] = replyParameters
	}

	return params
}
