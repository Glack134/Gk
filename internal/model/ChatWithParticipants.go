package model

type ChatWithParticipants struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	Participants []string `json:"participants"`
}
