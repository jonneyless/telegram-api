package models

import (
	"regexp"
	"strconv"
	"strings"
)

// InviteLink struct
type InviteLink struct {
	InviteLink              string `json:"invite_link"`
	Name                    string `json:"name"`
	Creator                 From   `json:"creator"`
	PendingJoinRequestCount int64  `json:"pending_join_request_count"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
}

func (l *InviteLink) GetCreatorId() int64 {
	re := regexp.MustCompile("user_(\\d+)")
	matched := re.MatchString(l.Name)
	if matched {
		str := strings.Replace(l.Name, "user_", "", -1)
		id, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}

		return id
	}

	return 0
}
