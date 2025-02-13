package service

import "github.com/polyk005/message/internal/api/repository"

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) SendNotification(userID int, message string) error {
	return s.repo.SendNotification(userID, message)
}

func (s *NotificationService) GetNotifications(userID string) ([]interface{}, error) {
	return s.repo.GetNotifications(userID)
}
