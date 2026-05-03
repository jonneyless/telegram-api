package requests

// Message struct
type Message struct {
	ChatId     int64   `json:"chat_id"`
	MessageId  int64   `json:"message_id"`
	MessageIds []int64 `json:"message_ids"`
}

func (p *Message) GetParams() map[string]interface{} {
	params := map[string]interface{}{
		"chat_id": p.ChatId,
	}

	messageIds := make([]int64, 0)

	if p.MessageIds != nil && len(p.MessageIds) > 0 {
		messageIds = p.MessageIds
	}

	if p.MessageId > 0 {
		messageIds = append(messageIds, p.MessageId)
	}

	params["message_ids"] = messageIds

	return params
}
