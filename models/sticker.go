package models

// Sticker struct
type Sticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Emoji        string `json:"emoji"`
	SetName      string `json:"set_name"`
	IsAnimated   bool   `json:"is_animated"`
	IsVideo      bool   `json:"is_video"`
	Type         string `json:"type"`
	Thumbnail    *Photo `json:"thumbnail"`
	Thumb        *Photo `json:"thumb"`
}
