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

func (s *ChatService) CreateChat(chatName string, userIDs ...int) (int, error) {
	if len(userIDs) == 2 {
		existingChatID, err := s.FindExistingChat(userIDs)
		if err != nil {
			return 0, err
		}
		if existingChatID != 0 {
			return existingChatID, nil
		}
	}

	return s.repo.CreateChat(chatName, userIDs...)
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

func (s *ChatService) ChatExistsBetweenUsers(userID1, userID2 int) (int, error) {
	return s.repo.ChatExists(userID1)
}

func (s *ChatService) GetUserIDByUsername(username string) (int, error) {
	return s.repo.GetUserIDByUsername(username)
}

func (s *ChatService) FindExistingChat(userIDs []int) (int, error) {
	return s.repo.FindExistingChat(userIDs)
}

func (s *ChatService) DeleteChatForAll(chatID int) error {
	return s.repo.DeleteChatForAll(chatID)
}

func (s *ChatService) DeleteChatForUser(chatID, userID int) error {
	return s.repo.DeleteChatForUser(chatID, userID)
}
