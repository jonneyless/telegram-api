package requests

import (
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type EditMessage struct {
	ChatId      int64                      `json:"chat_id"`
	MessageId   int64                      `json:"message_id"`
	Text        string                     `json:"text"`
	ParseMode   *string                    `json:"parse_mode"`
	Buttons     [][]map[string]interface{} `json:"buttons"`
	ReplyMarkup *models.ReplyMarkup        `json:"reply_markup"`
}

func (p *EditMessage) GetParams() map[string]interface{} {
	parseMode := "html"
	if p.ParseMode != nil {
		parseMode = *p.ParseMode
	}

	params := map[string]interface{}{
		"chat_id":    p.ChatId,
		"message_id": p.MessageId,
		"text":       p.Text,
		"parse_mode": parseMode,
	}

	if p.ReplyMarkup != nil {
		params["reply_markup"] = utils.Struct2Map(p.ReplyMarkup)
	} else if len(p.Buttons) > 0 {
		params["reply_markup"] = map[string]interface{}{
			"inline_keyboard": p.Buttons,
		}
	}

	return params
}
