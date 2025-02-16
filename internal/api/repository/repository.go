package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
	CreateUser(country, email, username, password string) error
	GetUserByPhone(phone string) (string, error)
	GetUserByEmail(email string) (string, error)
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

type Subscription interface {
	CreateSubscription(userID int, plan string) (int, error)
	GetSubscription(userID int) (map[string]interface{}, error)
	CancelSubscription(subscriptionID int) error
}

type Repository struct {
	Authorization Authorization
	Chat          *ChatRepository
	Message       *MessageRepository
	Notification  *NotificationRepository
	Payment       *PaymentRepository
	Subscription  *SubscriptionRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Chat:          NewChatRepository(db.DB),
		Message:       NewMessageRepository(db.DB),
		Notification:  NewNotificationRepository(db.DB),
		Payment:       NewPaymentRepository(db.DB),
		Subscription:  NewSubscriptionRepository(db.DB),
	}
}
