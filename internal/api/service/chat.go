package service

import (
	"github.com/polyk005/message/internal/api/repository"
)

type ChatService struct {
	repo *repository.ChatRepository
}

func NewChatService(repo *repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChat(userID, participantID int, chatName string) (int, error) {
	return s.repo.CreateChat(userID, participantID, chatName)
}

func (s *ChatService) AddParticipant(chatID, userID int) error {
	return s.repo.AddParticipant(chatID, userID)
}

func (s *ChatService) GetChatsForUser(userID int) ([]map[string]interface{}, error) {
	chats, err := s.repo.GetUserChats(userID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, chat := range chats {
		result = append(result, map[string]interface{}{
			"id":   chat.ID,
			"name": chat.Name,
		})
	}

	return result, nil
}

func (s *ChatService) UserExists(username string) (bool, error) {
	return s.repo.UserExists(username)
}

func (s *ChatService) ChatExists(userID int) (int, error) {
	return s.repo.ChatExists(userID)
}

func (s *ChatService) GetUserIDByUsername(username string) (int, error) {
	return s.repo.GetUserIDByUsername(username)
}
