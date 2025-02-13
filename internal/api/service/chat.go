package service

import "github.com/polyk005/message/internal/api/repository"

type ChatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChat(name string, participantIDs []int) (int, error) {
	return s.repo.CreateChat(name, participantIDs)
}

func (s *ChatService) GetChat(chatID string) (interface{}, error) {
	return s.repo.GetChat(chatID)
}

func (s *ChatService) AddParticipant(chatID, participantID int) error {
	return s.repo.AddParticipant(chatID, participantID)
}
