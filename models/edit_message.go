package models

// EditedMessage struct
type EditedMessage struct {
	MessageID   int64        `json:"message_id"`
	From        From         `json:"from"`
	Chat        Chat         `json:"chat"`
	Date        int64        `json:"date"`
	EditDate    int64        `json:"edit_date"`
	Text        *string      `json:"text,omitempty"`
	SenderTag   *string      `json:"sender_tag,omitempty"`
	Caption     *string      `json:"caption,omitempty"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
	Entities    *[]Entities  `json:"entities,omitempty"`
	Photo       *[]Photo     `json:"photo,omitempty"`
	Sticker     *Sticker     `json:"sticker,omitempty"`
	Video       *Video       `json:"video,omitempty"`
	Audio       *Audio       `json:"audio,omitempty"`
	Voice       *Voice       `json:"voice,omitempty"`
	Document    *Document    `json:"document,omitempty"`
	Dice        *Dice        `json:"dice,omitempty"`
	Poll        *Poll        `json:"poll,omitempty"`
}

// IsMessage 普通消息
func (m *EditedMessage) IsMessage() bool {
	return m.IsTextMessage() || m.IsPhoto() || m.IsSticker() || m.IsDocument() || m.IsVideo() || m.IsAudio() || m.IsVoice()
}

// IsTextMessage 文本消息
func (m *EditedMessage) IsTextMessage() bool {
	return m.Text != nil
}

// IsPhoto 图片消息
func (m *EditedMessage) IsPhoto() bool {
	return m.Photo != nil
}

// IsSticker 贴图消息
func (m *EditedMessage) IsSticker() bool {
	return m.Sticker != nil
}

// IsVideo 视频消息
func (m *EditedMessage) IsVideo() bool {
	return m.Video != nil
}

// IsAudio 音频消息
func (m *EditedMessage) IsAudio() bool {
	return m.Audio != nil
}

// IsVoice 语音消息
func (m *EditedMessage) IsVoice() bool {
	return m.Voice != nil
}

// IsDocument 文档消息
func (m *EditedMessage) IsDocument() bool {
	return m.Document != nil
}

func (m *EditedMessage) ChatId() int64 {
	return m.Chat.ID
}

func (m *EditedMessage) FromId() int64 {
	return m.From.ID
}

func (m *EditedMessage) TextString() string {
	if m.Text != nil {
		return *m.Text
	}

	if m.Caption != nil {
		return *m.Caption
	}

	return ""
}

func (m *EditedMessage) FileInfo() *FileInfo {
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
