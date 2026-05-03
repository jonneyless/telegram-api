package models

// ReceiveMessage struct
type ReceiveMessage struct {
	UpdateID        int64            `json:"update_id"`
	Message         *Message         `json:"message"`
	ChatMember      *ChatMember      `json:"chat_member"`
	MyChatMember    *ChatMember      `json:"my_chat_member"`
	EditedMessage   *EditedMessage   `json:"edited_message"`
	CallbackQuery   *CallbackQuery   `json:"callback_query"`
	ChatJoinRequest *ChatJoinRequest `json:"chat_join_request"`
	ChannelPost     *ChannelPost     `json:"channel_post"`
	Poll            *Poll            `json:"poll"`
	PollAnswer      *PollAnswer      `json:"poll_answer"`
}

func (m *ReceiveMessage) IsMessage() bool {
	return m.Message != nil
}

func (m *ReceiveMessage) IsChatMember() bool {
	return m.ChatMember != nil
}

func (m *ReceiveMessage) IsMyChatMember() bool {
	return m.MyChatMember != nil
}

func (m *ReceiveMessage) IsCallbackQuery() bool {
	return m.CallbackQuery != nil
}

func (m *ReceiveMessage) IsEditedMessage() bool {
	return m.EditedMessage != nil
}

func (m *ReceiveMessage) IsChatJoinRequest() bool {
	return m.ChatJoinRequest != nil
}

func (m *ReceiveMessage) IsPoll() bool {
	return m.Poll != nil
}

func (m *ReceiveMessage) IsPollAnswer() bool {
	return m.PollAnswer != nil
}
