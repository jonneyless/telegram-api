package models

// ChatJoinRequest struct
type ChatJoinRequest struct {
	From       From       `json:"from"`
	Chat       Chat       `json:"chat"`
	Date       int64      `json:"date"`
	InviteLink InviteLink `json:"invite_link"`
}
