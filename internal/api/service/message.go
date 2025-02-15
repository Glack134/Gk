package service

import "github.com/polyk005/message/internal/api/repository"

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SendMessage(chatID, userID int, content string) (int, error) {
	return s.repo.SendMessage(chatID, userID, content)
}

func (s *MessageService) EditMessage(messageID int, content string) error {
	return s.repo.EditMessage(messageID, content)
}

func (s *MessageService) DeleteMessage(messageID string) error {
	return s.repo.DeleteMessage(messageID)
}
