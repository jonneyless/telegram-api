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

func (m *ReceiveMessage) Text() string {
	if m.Message.IsMessage() {
		if m.Message.Text != nil {
			return *m.Message.Text
		}

		if m.Message.Caption != nil {
			return *m.Message.Caption
		}
	}

	if m.IsEditedMessage() {
		return m.EditedMessage.Text
	}

	return ""
}

func (m *ReceiveMessage) From() *From {
	if m.IsMessage() {
		return &m.Message.From
	}

	if m.IsEditedMessage() {
		return &m.EditedMessage.From
	}

	if m.IsCallbackQuery() {
		return &m.CallbackQuery.Message.From
	}

	if m.IsChatJoinRequest() {
		return &m.ChatJoinRequest.From
	}

	if m.IsChatMember() {
		return &m.ChatMember.From
	}

	if m.IsMyChatMember() {
		return &m.MyChatMember.From
	}

	return nil
}

func (m *ReceiveMessage) FromId() int64 {
	from := m.From()
	if from != nil {
		return from.ID
	}

	return 0
}

func (m *ReceiveMessage) Chat() *Chat {
	if m.IsMessage() {
		return &m.Message.Chat
	}

	if m.IsEditedMessage() {
		return &m.EditedMessage.Chat
	}

	if m.IsCallbackQuery() {
		return &m.CallbackQuery.Message.Chat
	}

	if m.IsChatJoinRequest() {
		return &m.ChatJoinRequest.Chat
	}

	if m.IsChatMember() {
		return &m.ChatMember.Chat
	}

	if m.IsMyChatMember() {
		return &m.MyChatMember.Chat
	}

	return nil
}

func (m *ReceiveMessage) ChatId() int64 {
	chat := m.Chat()
	if chat != nil {
		return chat.ID
	}

	return 0
}
