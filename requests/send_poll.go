package requests

import (
	"telegram_api/models"
	"telegram_api/utils"
)

type SendPoll struct {
	ChatId                int64               `json:"chat_id"`
	Question              string              `json:"question"`
	Options               []InputPollOption   `json:"options"`
	Type                  *string             `json:"type"`
	OpenPeriod            *int64              `json:"open_period"`
	CloseDate             *int64              `json:"close_date"`
	IsAnonymous           *bool               `json:"is_anonymous"`
	IsClosed              *bool               `json:"is_closed"`
	AllowsMultipleAnswers *bool               `json:"allows_multiple_answers"`
	Explanation           *string             `json:"explanation"`
	DisableNotification   *bool               `json:"disable_notification"`
	ProtectContent        *bool               `json:"protect_content"`
	ReplyParameters       *ReplyParameters    `json:"reply_parameters"`
	ReplyMarkup           *models.ReplyMarkup `json:"reply_markup"`
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
