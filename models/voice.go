package models

// Voice struct
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	Duration     int64  `json:"duration"`
	MimeType     string `json:"mime_type"`
}
