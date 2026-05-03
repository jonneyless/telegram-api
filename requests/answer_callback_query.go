package requests

type AnswerCallbackQuery struct {
	CallbackQueryId string  `json:"callback_query_id"`
	Text            string  `json:"text"`
	ShowAlert       bool    `json:"show_alert"`
	Url             *string `json:"url"`
	CacheTime       *int64  `json:"cache_time"`
}

func (p *AnswerCallbackQuery) GetParams() map[string]interface{} {
	params := map[string]interface{}{
		"callback_query_id": p.CallbackQueryId,
		"text":              p.Text,
		"show_alert":        p.ShowAlert,
	}

	if p.Url != nil {
		params["url"] = p.Url
	}

	if p.CacheTime != nil {
		params["cache_time"] = p.CacheTime
	}

	return params
}
