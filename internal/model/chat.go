package model

type Chat struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Participants []User `json:"participants"`
}
