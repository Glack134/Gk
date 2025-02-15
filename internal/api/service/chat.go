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

func (s *ChatService) CreateChat(name string) (int, error) {
	return s.repo.CreateChat(name)
}

func (s *ChatService) AddParticipant(chatID, userID int) error {
	return s.repo.AddParticipant(chatID, userID)
}

func (s *ChatService) GetChatsForUser(userID int) ([]map[string]interface{}, error) {
	return s.repo.GetChatsForUser(userID)
}
