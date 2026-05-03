package models

// Video struct
type Video struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileName     string `json:"file_name"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Duration     int64  `json:"duration"`
	MimeType     string `json:"mime_type"`
	Thumbnail    *Photo `json:"thumbnail"`
	Thumb        *Photo `json:"thumb"`
}
