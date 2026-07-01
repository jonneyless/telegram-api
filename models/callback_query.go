package models

import (
	"regexp"
	"strings"
)

// CallbackQuery struct
type CallbackQuery struct {
	Id              string   `json:"id"`
	From            From     `json:"from"`
	Data            string   `json:"data"`
	Message         *Message `json:"message,omitempty"`
	ChatInstance    *string  `json:"chat_instance,omitempty"`
	InlineMessageId *string  `json:"inline_message_id,omitempty"`
}

func (q *CallbackQuery) IsCallback(callback string) bool {
	return q.Data == callback
}

func (q *CallbackQuery) ChatId() int64 {
	return q.Message.Chat.ID
}

func (q *CallbackQuery) FromId() int64 {
	return q.From.ID
}

func (q *CallbackQuery) Contains(substr string) bool {
	return strings.Contains(q.Data, substr)
}

func (q *CallbackQuery) HasPrefix(prefix string) bool {
	return strings.HasPrefix(q.Data, prefix)
}

func (q *CallbackQuery) MatchQuery(key string, params ...bool) (bool, []string) {
	isPattern := false
	if len(params) > 0 {
		isPattern = params[0]
	}

	if isPattern {
		re := regexp.MustCompile(key)
		matches := re.FindStringSubmatch(q.Data)
		if matches != nil {
			return true, matches
		}
	} else {
		if q.Data == key {
			return true, nil
		}
	}

	return false, nil
}
