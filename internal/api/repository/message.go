package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) GetMessages(chatID string) ([]Message, error) {
	var messages []Message
	query := "SELECT * FROM messages WHERE chat_id = $1 ORDER BY created_at ASC"
	err := r.db.Select(&messages, query, chatID)
	return messages, err
}

func (r *MessageRepository) SendMessage(chatID, chatParticipantID int, content string) (int, error) {
	// Проверяем, существует ли chat_participant_id в таблице chat_participants
	var participantID int
	err := r.db.QueryRow("SELECT id FROM chat_participants WHERE id = $1", chatParticipantID).Scan(&participantID)
	if err != nil {
		return 0, fmt.Errorf("chat participant with id %d does not exist: %v", chatParticipantID, err)
	}

	// Если участник существует, вставляем сообщение
	var messageID int
	query := `INSERT INTO messages (chat_id, chat_participant_id, content) VALUES ($1, $2, $3) RETURNING id`
	err = r.db.QueryRow(query, chatID, chatParticipantID, content).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %v", err)
	}

	return messageID, nil
}

func (r *MessageRepository) EditMessage(messageID int, content string) error {
	query := `UPDATE messages SET content = $1 WHERE id = $2`
	_, err := r.db.Exec(query, content, messageID)
	return err
}

func (r *MessageRepository) DeleteMessage(messageID string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, messageID)
	return err
}
