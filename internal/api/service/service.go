package service

import "github.com/polyk005/message/internal/api/repository"

type Authorization interface {
	SendVerificationCode(identifier string) error
	VerifyCode(identifier, code string) bool
	SignUp(identifier, code string) error
	SignIn(identifier, code string) error
}

type Chat interface {
	CreateChat(name string, participantIDs []int) (int, error)
	GetChat(chatID string) (interface{}, error)
	AddParticipant(chatID, participantID int) error
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

type Service struct {
	Authorization
	Chat
	Message
	Notification
	Payment
}

func NewService(repos *repository.Repository, sms *SMSService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, sms),
		Chat:          NewChatService(repository.ChatRepository{}),
		Message:       NewMessageService(repository.MessageRepository{}),
		Notification:  NewNotificationService(repository.NotificationRepository{}),
		Payment:       NewPaymentService(repository.PaymentRepository{}),
	}
}
