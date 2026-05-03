package models

import (
	"regexp"
	"strings"
)

// Message struct
type Message struct {
	MessageID            int64           `json:"message_id"`
	From                 From            `json:"from"`
	Chat                 Chat            `json:"chat"`
	Date                 int64           `json:"date"`
	Text                 *string         `json:"text"`
	Caption              *string         `json:"caption"`
	ReplyMarkup          *ReplyMarkup    `json:"reply_markup"`
	ReplyToMessage       *ReplyToMessage `json:"reply_to_message"`
	PinnedMessage        *Message        `json:"pinned_message"`
	ForwardOrigin        *Origin         `json:"forward_origin"`
	ForwardFrom          *From           `json:"forward_from"`
	ForwardFromChat      *Chat           `json:"forward_from_chat"`
	ForwardFromMessageID *int64          `json:"forward_from_message_id"`
	ForwardDate          *int64          `json:"forward_date"`
	ExternalReply        *ExternalReply  `json:"external_reply"`
	Quote                *Quote          `json:"quote"`
	LeftChatParticipant  *From           `json:"left_chat_participant"`
	LeftChatMember       *From           `json:"left_chat_member"`
	NewChatParticipant   *From           `json:"new_chat_participant"`
	NewChatMember        *From           `json:"new_chat_member"`
	NewChatMembers       *[]From         `json:"new_chat_members"`
	NewChatTitle         *string         `json:"new_chat_title"`
	NewChatPhoto         *[]Photo        `json:"new_chat_photo"`
	Entities             *[]Entities     `json:"entities"`
	Photo                *[]Photo        `json:"photo"`
	Sticker              *Sticker        `json:"sticker"`
	Video                *Video          `json:"video"`
	Audio                *Audio          `json:"audio"`
	Voice                *Voice          `json:"voice"`
	Document             *Document       `json:"document"`
	Dice                 *Dice           `json:"dice"`
	Poll                 *Poll           `json:"poll"`
}

// ReplyToMessage struct
type ReplyToMessage struct {
	MessageID   int64        `json:"message_id"`
	From        From         `json:"from"`
	Chat        Chat         `json:"chat"`
	Date        int64        `json:"date"`
	Text        *string      `json:"text"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup"`
	Entities    *[]Entities  `json:"entities"`
}

// ExternalReply struct
type ExternalReply struct {
	MessageID int64  `json:"message_id"`
	Origin    Origin `json:"origin"`
	Chat      Chat   `json:"chat"`
}

// Origin struct
type Origin struct {
	MessageID  int64  `json:"message_id"`
	Type       string `json:"type"`
	SenderUser From   `json:"sender_user"`
	Chat       Chat   `json:"chat"`
	Date       int64  `json:"date"`
}

// Quote struct
type Quote struct {
	Text     string `json:"text"`
	Position int64  `json:"position"`
	IsManual bool   `json:"is_manual"`
}

// Entities struct
type Entities struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
}

// IsCommand 匹配命令
func (m *Message) IsCommand(command string) bool {
	if m.Text == nil {
		return false
	}

	text := strings.TrimLeft(*m.Text, "/")

	return text == command
}

func (m *Message) Contains(substr string) bool {
	return strings.Contains(*m.Text, substr)
}

func (m *Message) HasPrefix(prefix string) bool {
	return strings.HasPrefix(*m.Text, prefix)
}

// IsGroupMange 群管命令
func (m *Message) IsGroupMange(key string, params ...bool) (bool, []string) {
	isPattern := false
	if len(params) > 0 {
		isPattern = params[0]
	}

	if isPattern {
		re := regexp.MustCompile(key)
		matches := re.FindStringSubmatch(*m.Text)
		if matches != nil {
			return true, matches
		}
	} else {
		if *m.Text == key {
			return true, nil
		}
	}

	return false, nil
}

// IsPrivate 私聊
func (m *Message) IsPrivate() bool {
	return m.Chat.Type == "private"
}

// IsReplyToMessage 回复消息
func (m *Message) IsReplyToMessage() bool {
	return m.ReplyToMessage != nil
}

// IsActionMessage 系统消息
func (m *Message) IsActionMessage() bool {
	return m.IsNewChatMember() || m.IsLeftChatMember() || m.IsUpdateChatTitle() || m.IsUpdateChatPhoto() || m.IsPinnedMessage()
}

// IsNewChatMember 用户入群
func (m *Message) IsNewChatMember() bool {
	return m.NewChatParticipant != nil
}

// IsLeftChatMember 用户离群
func (m *Message) IsLeftChatMember() bool {
	return m.LeftChatParticipant != nil
}

// IsUpdateChatTitle 更改群名
func (m *Message) IsUpdateChatTitle() bool {
	return m.NewChatTitle != nil
}

// IsUpdateChatPhoto 更改群头像
func (m *Message) IsUpdateChatPhoto() bool {
	return m.NewChatPhoto != nil || len(*m.NewChatPhoto) > 0
}

// IsPinnedMessage 消息置顶
func (m *Message) IsPinnedMessage() bool {
	return m.PinnedMessage != nil
}

// IsMessage 普通消息
func (m *Message) IsMessage() bool {
	return m.IsTextMessage() || m.IsPhoto() || m.IsSticker() || m.IsDocument() || m.IsVideo() || m.IsAudio() || m.IsVoice()
}

// IsTextMessage 文本消息
func (m *Message) IsTextMessage() bool {
	return m.Text != nil
}

// IsPhoto 图片消息
func (m *Message) IsPhoto() bool {
	return m.Photo != nil
}

// IsSticker 贴图消息
func (m *Message) IsSticker() bool {
	return m.Sticker != nil
}

// IsVideo 视频消息
func (m *Message) IsVideo() bool {
	return m.Video != nil
}

// IsAudio 音频消息
func (m *Message) IsAudio() bool {
	return m.Audio != nil
}

// IsVoice 语音消息
func (m *Message) IsVoice() bool {
	return m.Voice != nil
}

// IsDocument 文档消息
func (m *Message) IsDocument() bool {
	return m.Document != nil
}

func (m *Message) ChatId() int64 {
	return m.Chat.ID
}

func (m *Message) FromId() int64 {
	return m.From.ID
}
