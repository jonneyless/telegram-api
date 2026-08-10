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
	Text                 *string         `json:"text,omitempty"`
	Caption              *string         `json:"caption,omitempty"`
	ReplyMarkup          *ReplyMarkup    `json:"reply_markup,omitempty"`
	ReplyToMessage       *ReplyToMessage `json:"reply_to_message,omitempty"`
	PinnedMessage        *Message        `json:"pinned_message,omitempty"`
	ForwardOrigin        *Origin         `json:"forward_origin,omitempty"`
	ForwardFrom          *From           `json:"forward_from,omitempty"`
	ForwardFromChat      *Chat           `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID *int64          `json:"forward_from_message_id,omitempty"`
	ForwardDate          *int64          `json:"forward_date,omitempty"`
	ExternalReply        *ExternalReply  `json:"external_reply,omitempty"`
	Quote                *Quote          `json:"quote,omitempty"`
	LeftChatParticipant  *From           `json:"left_chat_participant,omitempty"`
	LeftChatMember       *From           `json:"left_chat_member,omitempty"`
	NewChatParticipant   *From           `json:"new_chat_participant,omitempty"`
	NewChatMember        *From           `json:"new_chat_member,omitempty"`
	NewChatMembers       *[]From         `json:"new_chat_members,omitempty"`
	NewChatTitle         *string         `json:"new_chat_title,omitempty"`
	NewChatPhoto         *[]Photo        `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto      *bool           `json:"delete_chat_photo,omitempty"`
	Entities             *[]Entities     `json:"entities,omitempty"`
	Photo                *[]Photo        `json:"photo,omitempty"`
	Sticker              *Sticker        `json:"sticker,omitempty"`
	Video                *Video          `json:"video,omitempty"`
	Audio                *Audio          `json:"audio,omitempty"`
	Voice                *Voice          `json:"voice,omitempty"`
	Document             *Document       `json:"document,omitempty"`
	Dice                 *Dice           `json:"dice,omitempty"`
	Poll                 *Poll           `json:"poll,omitempty"`
	ChatShared           *ChatShared     `json:"chat_shared,omitempty"`
}

// ReplyToMessage struct
type ReplyToMessage struct {
	MessageID   int64        `json:"message_id"`
	From        From         `json:"from"`
	Chat        Chat         `json:"chat"`
	Date        int64        `json:"date"`
	Text        *string      `json:"text,omitempty"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
	Entities    *[]Entities  `json:"entities,omitempty"`
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
	User   *From  `json:"user,omitempty"`
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
	return m.Chat.IsPrivate()
}

// IsReplyToMessage 回复消息
func (m *Message) IsReplyToMessage() bool {
	return m.ReplyToMessage != nil
}

// IsForwardMessage 转发消息
func (m *Message) IsForwardMessage() bool {
	return m.ForwardOrigin != nil
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
	return m.NewChatPhoto != nil || m.DeleteChatPhoto != nil
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

func (m *Message) TextString() string {
	if m.Text != nil {
		return *m.Text
	}

	if m.Caption != nil {
		return *m.Caption
	}

	return ""
}

type FileInfo struct {
	Type         string  `json:"type"`
	FileId       string  `json:"file_id"`
	FileUniqueId string  `json:"file_unique_id"`
	FileSize     int64   `json:"file_size"`
	FileName     *string `json:"file_name"`
	Width        *int64  `json:"width"`
	Height       *int64  `json:"height"`
	MimeType     *string `json:"mime_type"`
	Duration     *int64  `json:"duration"`
}

func (m *Message) FileInfo() *FileInfo {
	info := &FileInfo{}

	if m.IsPhoto() {
		info.Type = "photo"
		for _, photo := range *m.Photo {
			if info.FileSize < photo.FileSize {
				info.FileId = photo.FileID
				info.FileUniqueId = photo.FileUniqueID
				info.FileSize = photo.FileSize
				info.Width = &photo.Width
				info.Height = &photo.Height
			}
		}
	} else if m.IsVideo() {
		info.Type = "video"
		info.FileId = m.Video.FileID
		info.FileUniqueId = m.Video.FileUniqueID
		info.FileSize = m.Video.FileSize
		info.Width = &m.Video.Width
		info.Height = &m.Video.Height
		info.FileName = &m.Video.FileName
		info.Duration = &m.Video.Duration
		info.MimeType = &m.Video.MimeType
	} else if m.IsAudio() {
		info.Type = "audio"
		info.FileId = m.Audio.FileID
		info.FileUniqueId = m.Audio.FileUniqueID
		info.FileSize = m.Audio.FileSize
		info.FileName = &m.Audio.FileName
		info.Duration = &m.Audio.Duration
		info.MimeType = &m.Audio.MimeType
	} else if m.IsSticker() {
		info.Type = "sticker"
		info.FileId = m.Sticker.FileID
		info.FileUniqueId = m.Sticker.FileUniqueID
		info.FileSize = m.Sticker.FileSize
		info.Width = &m.Sticker.Width
		info.Height = &m.Sticker.Height
	} else if m.IsDocument() {
		info.Type = "document"
		info.FileId = m.Document.FileID
		info.FileUniqueId = m.Document.FileUniqueID
		info.FileSize = m.Document.FileSize
		info.FileName = &m.Document.FileName
		info.MimeType = &m.Document.MimeType
	} else if m.IsVoice() {
		info.Type = "voice"
		info.FileId = m.Voice.FileID
		info.FileUniqueId = m.Voice.FileUniqueID
		info.FileSize = m.Voice.FileSize
		info.Duration = &m.Voice.Duration
		info.MimeType = &m.Voice.MimeType
	} else {
		return nil
	}

	return info
}
