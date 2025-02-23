package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/message/internal/model"
)

type User interface {
	GetUserID(userID int) (*model.User, error)
	UpdateUser(user *model.User) error
	UpdateUserEmail(userID int, newEmail string) error
	GetTokenResetPassword(email string) (int, time.Time, error)
}

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(email, password string) (model.User, error)
}
type Chat interface {
	CreateChat(chatName string, userIDs ...int) (int, error)
	AddUserToChat(chatID, userID int) error
	AddParticipant(chatID, userID int) error
	GetUserChats(userID int) ([]model.ChatWithParticipants, error)
	UserExists(username string) (bool, error)
	ChatExists(userID int) (int, error)
	GetUserIDByUsername(username string) (int, error)
	ChatExistsBetweenUsers(userIDs ...int) (int, error)
	DeleteChatForAll(chatID int) error
	DeleteChatForUser(chatID, userID int) error
	FindExistingChat(userID []int) (int, error)
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
	User          *UserRepository
	Chat          *ChatRepository
	Message       *MessageRepository
	Notification  *NotificationRepository
	Payment       *PaymentRepository
	Subscription  *SubscriptionRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:          NewUserRepository(db.DB),
		Authorization: NewAuthPostgres(db),
		Chat:          NewChatRepository(db.DB),
		Message:       NewMessageRepository(db.DB),
		Notification:  NewNotificationRepository(db.DB),
		Payment:       NewPaymentRepository(db.DB),
		Subscription:  NewSubscriptionRepository(db.DB),
	}
}
