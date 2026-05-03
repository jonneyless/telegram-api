package requests

import (
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type SendMessage struct {
	ChatId              int64                      `json:"chat_id"`
	Text                string                     `json:"text"`
	Photo               string                     `json:"photo"`
	Video               string                     `json:"video"`
	Audio               string                     `json:"audio"`
	Document            string                     `json:"document"`
	ParseMode           *string                    `json:"parse_mode"`
	ReplyMarkup         *models.ReplyMarkup        `json:"reply_markup"`
	Buttons             [][]map[string]interface{} `json:"buttons"`
	DisableNotification bool                       `json:"disable_notification"`
	LinkPreviewOptions  *LinkPreviewOptions        `json:"link_preview_options"`
	ReplyParameters     *ReplyParameters           `json:"reply_parameters"`
}

func (p *SendMessage) GetParams() map[string]interface{} {
	parseMode := "html"
	if p.ParseMode != nil {
		parseMode = *p.ParseMode
	}

	params := map[string]interface{}{
		"chat_id":              p.ChatId,
		"text":                 p.Text,
		"parse_mode":           parseMode,
		"disable_notification": p.DisableNotification,
	}

	if p.Photo != "" || p.Video != "" || p.Audio != "" || p.Document != "" {
		params = map[string]interface{}{
			"chat_id":              p.ChatId,
			"caption":              p.Text,
			"parse_mode":           parseMode,
			"disable_notification": p.DisableNotification,
		}

		if p.Photo != "" {
			params["photo"] = p.Photo
		} else if p.Video != "" {
			params["video"] = p.Video
		} else if p.Audio != "" {
			params["audio"] = p.Audio
		} else if p.Document != "" {
			params["document"] = p.Document
		}
	}

	if p.ReplyMarkup != nil {
		params["reply_markup"] = utils.Struct2Map(p.ReplyMarkup)
	} else if len(p.Buttons) > 0 {
		params["reply_markup"] = map[string]interface{}{
			"inline_keyboard": p.Buttons,
		}
	}

	if p.LinkPreviewOptions != nil {
		linkPreviewOptions := map[string]interface{}{
			"is_disabled": p.LinkPreviewOptions.IsDisabled,
		}
		if p.LinkPreviewOptions.Url != nil {
			linkPreviewOptions["url"] = p.LinkPreviewOptions.Url
		}

		params["link_preview_options"] = linkPreviewOptions
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

type LinkPreviewOptions struct {
	IsDisabled bool    `json:"is_disabled"`
	Url        *string `json:"url"`
}

type ReplyParameters struct {
	MessageId int64  `json:"message_id"`
	ChatId    *int64 `json:"chat_id"`
}
