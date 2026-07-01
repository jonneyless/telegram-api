package requests

import (
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type SendPoll struct {
	ChatId                int64               `json:"chat_id"`
	Question              string              `json:"question"`
	Options               []InputPollOption   `json:"options"`
	Type                  *string             `json:"type,omitempty"`
	OpenPeriod            *int64              `json:"open_period,omitempty"`
	CloseDate             *int64              `json:"close_date,omitempty"`
	IsAnonymous           *bool               `json:"is_anonymous,omitempty"`
	IsClosed              *bool               `json:"is_closed,omitempty"`
	AllowsMultipleAnswers *bool               `json:"allows_multiple_answers,omitempty"`
	Explanation           *string             `json:"explanation,omitempty"`
	DisableNotification   *bool               `json:"disable_notification,omitempty"`
	ProtectContent        *bool               `json:"protect_content,omitempty"`
	ReplyParameters       *ReplyParameters    `json:"reply_parameters,omitempty"`
	ReplyMarkup           *models.ReplyMarkup `json:"reply_markup,omitempty"`
}

func (p *SendPoll) GetParams() map[string]interface{} {
	params := map[string]interface{}{
		"chat_id":                 p.ChatId,
		"question":                p.Question,
		"question_parse_mode":     "html",
		"type":                    "regular",
		"is_anonymous":            p.IsAnonymous,
		"is_closed":               p.IsClosed,
		"allows_multiple_answers": p.AllowsMultipleAnswers,
		"disable_notification":    p.DisableNotification,
		"protect_content":         p.ProtectContent,
	}

	if p.Type != nil {
		params["type"] = p.Type
	}

	if p.CloseDate != nil {
		params["close_date"] = p.CloseDate
	}

	if p.Explanation != nil {
		params["explanation"] = p.Explanation
	}

	options := make([]map[string]interface{}, 0)
	for _, option := range p.Options {
		item := utils.Struct2Map(option)
		options = append(options, item)
	}
	params["options"] = options

	return params
}

type InputPollOption struct {
	Text string `json:"text"`
}
