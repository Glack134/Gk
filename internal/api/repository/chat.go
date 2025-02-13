package repository

import (
	"database/sql"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(name string, participantIDs []int) (int, error) {
	return 1, nil
}

func (r *ChatRepository) GetChat(chatID string) (interface{}, error) {
	return nil, nil
}

func (r *ChatRepository) AddParticipant(chatID, participantID int) error {
	return nil
}
