package requests

import (
	"telegram_api/models"
	"telegram_api/utils"
)

type SendDice struct {
	ChatId              int64               `json:"chat_id"`
	Emoji               *string             `json:"emoji"`
	DisableNotification *bool               `json:"disable_notification"`
	ProtectContent      *bool               `json:"protect_content"`
	ReplyParameters     *ReplyParameters    `json:"reply_parameters"`
	ReplyMarkup         *models.ReplyMarkup `json:"reply_markup"`
}

func (p *SendDice) GetParams() map[string]interface{} {
	params := map[string]interface{}{
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
		replyParameters := map[string]interface{}{
			"message_id": p.ReplyParameters.MessageId,
		}
		if p.ReplyParameters.ChatId != nil {
			replyParameters["chat_id"] = p.ReplyParameters.ChatId
		}

		params["reply_parameters"] = replyParameters
	}

	return params
}
