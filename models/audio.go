package models

// Audio struct
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileName     string `json:"file_name"`
	Title        string `json:"title"`
	Duration     int64  `json:"duration"`
	Performer    string `json:"performer"`
	MimeType     string `json:"mime_type"`
	Thumbnail    *Photo `json:"thumbnail"`
}
