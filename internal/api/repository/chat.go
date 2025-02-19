package repository

import (
	"database/sql"

	"github.com/polyk005/message/internal/model"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(userID, participantID int, chatName string) (int, error) {
	// Создаем новый чат
	var newChatID int
	query := `INSERT INTO chats (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRow(query, chatName).Scan(&newChatID)
	if err != nil {
		return 0, err
	}

	err = r.AddUserToChat(newChatID, userID)
	if err != nil {
		return 0, err
	}

	err = r.AddUserToChat(newChatID, participantID)
	if err != nil {
		return 0, err
	}

	return newChatID, nil
}

func (r *ChatRepository) AddUserToChat(chatID, userID int) error {
	query := `INSERT INTO chat_participants (chat_id, user_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, chatID, userID)
	return err
}

func (r *ChatRepository) AddParticipant(chatID, userID int) error {
	return r.AddUserToChat(chatID, userID)
}

func (r *ChatRepository) GetUserChats(userID int) ([]model.Chat, error) {
	var chats []model.Chat
	query := `
		SELECT c.id, c.name 
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(&chat.ID, &chat.Name); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (r *ChatRepository) UserExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (r *ChatRepository) ChatExists(userID int) (int, error) {
	var chatID int
	query := `SELECT c.id 
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&chatID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func (r *ChatRepository) GetUserIDByUsername(username string) (int, error) {
	var userID int
	query := `SELECT id FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func (r *ChatRepository) ChatExistsBetweenUsers(userID1, userID2 int) (int, error) {
	var chatID int
	query := `
		SELECT c.id 
		FROM chats c
		JOIN chat_participants cp1 ON c.id = cp1.chat_id
		JOIN chat_participants cp2 ON c.id = cp2.chat_id
		WHERE cp1.user_id = $1 AND cp2.user_id = $2
	`
	err := r.db.QueryRow(query, userID1, userID2).Scan(&chatID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
