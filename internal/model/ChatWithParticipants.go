package model

type ChatWithParticipants struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}
