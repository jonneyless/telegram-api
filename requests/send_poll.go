package requests

import (
	"github.com/bytedance/sonic"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/utils"
)

type SendPoll struct {
	ChatId                int64               `json:"chat_id"`
	Question              string              `json:"question"`
	QuestionParseMode     string              `json:"question_parse_mode"`
	Options               []InputPollOption   `json:"options"`
	Type                  string              `json:"type,omitempty"`
	OpenPeriod            int64               `json:"open_period,omitempty"`
	CloseDate             int64               `json:"close_date,omitempty"`
	IsAnonymous           bool                `json:"is_anonymous"`
	IsClosed              bool                `json:"is_closed"`
	AllowsMultipleAnswers bool                `json:"allows_multiple_answers"`
	Explanation           string              `json:"explanation,omitempty"`
	DisableNotification   bool                `json:"disable_notification"`
	ProtectContent        bool                `json:"protect_content"`
	ReplyParameters       *ReplyParameters    `json:"reply_parameters,omitempty"`
	ReplyMarkup           *models.ReplyMarkup `json:"reply_markup,omitempty"`
	Buttons               [][]map[string]any  `json:"-"`
	ButtonType            string              `json:"-"`
}

func (p *SendPoll) MarshalJSON() ([]byte, error) {
	data := *p

	if data.QuestionParseMode == "" {
		data.QuestionParseMode = "html"
	}

	if data.Type == "" {
		data.Type = "regular"
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

type InputPollOption struct {
	Text string `json:"text"`
}
