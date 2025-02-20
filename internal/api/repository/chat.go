package repository

import (
	"database/sql"
	"strings"

	"github.com/polyk005/message/internal/model"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(chatName string, userIDs ...int) (int, error) {
	var newChatID int
	query := `INSERT INTO chats (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRow(query, chatName).Scan(&newChatID)
	if err != nil {
		return 0, err
	}

	// Добавляем участников в чат
	for _, userID := range userIDs {
		err := r.AddUserToChat(newChatID, userID)
		if err != nil {
			return 0, err
		}
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

func (r *ChatRepository) ChatExistsBetweenUsers(userIDs ...int) (int, error) {
	var chatID int
	query := `
        SELECT c.id 
        FROM chats c
        JOIN chat_participants cp ON c.id = cp.chat_id
        WHERE cp.user_id IN (` + strings.Repeat("?,", len(userIDs)-1) + `?)
        GROUP BY c.id
        HAVING COUNT(DISTINCT cp.user_id) = ?
    `
	args := make([]interface{}, len(userIDs)+1)
	for i, id := range userIDs {
		args[i] = id
	}
	args[len(userIDs)] = len(userIDs)

	err := r.db.QueryRow(query, args...).Scan(&chatID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func (r *ChatRepository) DeleteChatForAll(chatID int) error {
	query := `DELETE FROM chats WHERE id = $1`
	_, err := r.db.Exec(query, chatID)
	return err
}

func (r *ChatRepository) DeleteChatForUser(chatID, userID int) error {
	query := `DELETE FROM chat_participants WHERE chat_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, chatID, userID)
	return err
}
