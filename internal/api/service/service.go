package service

import "github.com/polyk005/message/internal/api/repository"

type Authorization interface {
	SendVerificationCode(identifier string) error
	VerifyCode(identifier, code string) bool
	SignUp(identifier, code string) error
	SignIn(identifier, code string) error
	ValidateToken(tokenString string) (int, error)
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

func NewService(repos *repository.Repository, sms *SMSService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, sms),
		Chat:          NewChatService(repos.Chat),
		Message:       NewMessageService(repos.Message),
		Notification:  NewNotificationService(repos.Notification),
		Payment:       NewPaymentService(repos.Payment),
		Subscription:  NewSubscriptionService(repos.Subscription),
	}
}
