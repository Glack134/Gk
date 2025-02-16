package service

import (
	"github.com/polyk005/message/internal/api/repository"
	"github.com/polyk005/message/internal/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	generatePasswordHash(password string) string
	HashPassword(password string) (string, error)
}

type Chat interface {
	CreateChat(name string) (int, error)
	AddParticipant(chatID, userID int) error
	GetChatsForUser(userID int) ([]map[string]interface{}, error)
}

type Message interface {
	SendMessage(chatID, userID int, content string) (int, error)
	EditMessage(messageID int, content string) error
	DeleteMessage(messageID string) error
}

type Notification interface {
	SendNotification(userID int, message string) error
	GetNotifications(userID string) ([]interface{}, error)
}

type Payment interface {
	CreatePayment(userID int, amount float64, purpose string) (int, error)
	GetPaymentStatus(paymentID string) (string, error)
}

type Subscription interface {
	CreateSubscription(userID int, plan string) (int, error)
	GetSubscription(userID int) (map[string]interface{}, error)
	CancelSubscription(subscriptionID int) error
}

type Service struct {
	Authorization Authorization
	Chat          Chat
	Message       Message
	Notification  Notification
	Payment       Payment
	Subscription  Subscription
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Chat:          NewChatService(repos.Chat),
		Message:       NewMessageService(repos.Message),
		Notification:  NewNotificationService(repos.Notification),
		Payment:       NewPaymentService(repos.Payment),
		Subscription:  NewSubscriptionService(repos.Subscription),
	}
}
