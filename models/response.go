package models

type Response[T any] struct {
	Ok          bool   `json:"ok"`
	Result      *T     `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}
