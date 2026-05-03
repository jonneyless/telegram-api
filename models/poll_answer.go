package models

type PollAnswer struct {
	PollId    string  `json:"poll_id"`
	User      From    `json:"user"`
	OptionIds []int64 `json:"option_ids"`
}
