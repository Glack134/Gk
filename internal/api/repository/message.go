package repository

import (
	"database/sql"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) SendMessage(chatID, userID int, content string) (int, error) {
	return 1, nil
}

func (r *MessageRepository) EditMessage(messageID int, content string) error {
	return nil
}

func (r *MessageRepository) DeleteMessage(messageID string) error {
	return nil
}
