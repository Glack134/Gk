package repository

import (
	"database/sql"
	"time"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(name string) (int, error) {
	var chatID int
	query := `INSERT INTO chats (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRow(query, name).Scan(&chatID)
	return chatID, err
}

func (r *ChatRepository) AddParticipant(chatID, userID int) error {
	query := `INSERT INTO chat_participants (chat_id, user_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, chatID, userID)
	return err
}

func (r *ChatRepository) GetChatsForUser(userID int) ([]map[string]interface{}, error) {
	query := `
		SELECT c.id, c.name, c.created_at 
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []map[string]interface{}
	for rows.Next() {
		var (
			id        int
			name      string
			createdAt time.Time
		)
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			return nil, err
		}
		chats = append(chats, map[string]interface{}{
			"id":         id,
			"name":       name,
			"created_at": createdAt,
		})
	}

	return chats, nil
}
