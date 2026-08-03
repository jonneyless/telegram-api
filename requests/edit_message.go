package requests

import (
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type EditMessageText struct {
	ChatId             int64               `json:"chat_id"`
	MessageId          int64               `json:"message_id"`
	Text               string              `json:"text,omitempty"`
	ParseMode          string              `json:"parse_mode,omitempty"`
	Buttons            [][]map[string]any  `json:"buttons"`
	ReplyMarkup        *models.ReplyMarkup `json:"reply_markup,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

func (p *EditMessageText) GetParams() map[string]any {
	parseMode := "html"
	if p.ParseMode != "" {
		parseMode = p.ParseMode
	}

	params := map[string]any{
		"chat_id":    p.ChatId,
		"message_id": p.MessageId,
		"text":       p.Text,
		"parse_mode": parseMode,
	}

	if p.ReplyMarkup != nil {
		params["reply_markup"] = utils.Struct2Map(p.ReplyMarkup)
	} else if len(p.Buttons) > 0 {
		params["reply_markup"] = map[string]any{
			"inline_keyboard": p.Buttons,
		}
	}

	if p.LinkPreviewOptions != nil {
		linkPreviewOptions := map[string]any{
			"is_disabled": p.LinkPreviewOptions.IsDisabled,
		}
		if p.LinkPreviewOptions.Url != "" {
			linkPreviewOptions["url"] = p.LinkPreviewOptions.Url
		}

		params["link_preview_options"] = linkPreviewOptions
	}

	return params
}

type EditMessageCaption struct {
	ChatId             int64               `json:"chat_id"`
	MessageId          int64               `json:"message_id"`
	Caption            string              `json:"caption,omitempty"`
	ParseMode          string              `json:"parse_mode,omitempty"`
	ReplyMarkup        *models.ReplyMarkup `json:"reply_markup,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

type EditMessageMedia struct {
	ChatId      int64               `json:"chat_id"`
	MessageId   int64               `json:"message_id"`
	Media       InputMedia          `json:"media"`
	ReplyMarkup *models.ReplyMarkup `json:"reply_markup,omitempty"`
}

type InputMedia struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}
