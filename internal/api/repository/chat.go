package repository

import (
	"database/sql"
	"sort"

	"github.com/lib/pq"
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

func (r *ChatRepository) GetUserChats(userID int) ([]model.ChatWithParticipants, error) {
	var chats []model.ChatWithParticipants

	// Получаем имя пользователя
	var username string
	err := r.db.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&username)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT c.id, c.name, array_agg(DISTINCT u.username) as participants
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		JOIN users u ON cp.user_id = u.id
		GROUP BY c.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.ChatWithParticipants
		var participants pq.StringArray
		if err := rows.Scan(&chat.ID, &chat.Name, &participants); err != nil {
			return nil, err
		}

		// Добавляем реальное имя пользователя в список участников, если его там нет
		participants = append(participants, username)
		chat.Participants = participants

		// Удаляем дубликаты
		uniqueParticipants := make(map[string]struct{})
		for _, participant := range chat.Participants {
			uniqueParticipants[participant] = struct{}{}
		}

		// Преобразуем обратно в срез
		chat.Participants = make([]string, 0, len(uniqueParticipants))
		for participant := range uniqueParticipants {
			chat.Participants = append(chat.Participants, participant)
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
	 WHERE (cp1.user_id = $1 AND cp2.user_id = $2)
		OR (cp1.user_id = $2 AND cp2.user_id = $1)
	 LIMIT 1
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

func (r *ChatRepository) FindExistingChat(userIDs []int) (int, error) {
	if len(userIDs) != 2 {
		return 0, nil
	}

	sort.Ints(userIDs)
	query := `
    SELECT c.id
    FROM chats c
    JOIN chat_participants cp1 ON c.id = cp1.chat_id
    JOIN chat_participants cp2 ON c.id = cp2.chat_id
    WHERE cp1.user_id = $1 AND cp2.user_id = $2
    AND NOT EXISTS (
        SELECT 1
        FROM chat_participants cp3
        WHERE cp3.chat_id = c.id
        AND cp3.user_id NOT IN ($1, $2)
    )
    LIMIT 1;` // Добавлена закрывающая скобка и исправлен порядок LIMIT

	var chatID int
	err := r.db.QueryRow(query, userIDs[0], userIDs[1]).Scan(&chatID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
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
