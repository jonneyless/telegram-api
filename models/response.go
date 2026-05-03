package models

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiResponse struct {
	Ok     bool `json:"ok"`
	Result any  `json:"result"`
}

type ApiErrorResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

// MessageResponse struct
type MessageResponse struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}

// ChatInviteLinkResponse struct
type ChatInviteLinkResponse struct {
	Ok     bool       `json:"ok"`
	Result InviteLink `json:"result"`
}

// ChatResponse struct
type ChatResponse struct {
	Ok     bool `json:"ok"`
	Result Chat `json:"result"`
}

// ChatMemberResponse struct
type ChatMemberResponse struct {
	Ok     bool   `json:"ok"`
	Result Member `json:"result"`
}

// ChatMembersResponse struct
type ChatMembersResponse struct {
	Ok     bool     `json:"ok"`
	Result []Member `json:"result"`
}

// FileResponse struct
type FileResponse struct {
	Ok     bool `json:"ok"`
	Result File `json:"result"`
}
