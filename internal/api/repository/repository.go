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
}

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(email string, password string, checkPassword bool) (model.User, error)
	GetUserIDByToken(token string) (int, error)
	UpdatePasswordUserByID(userID int, newPasswordHash string) error
	MarkTokenAsUsed(token string) error
	IsTokenUsed(token string) (bool, error)
	GetLastSentTime(token string) (time.Time, error)
	UpdateTwoFASecret(userID int, secret string) error
	GetTwoFASecret(userID int) (string, error)
	DisableTwoFA(userID int) error
	IsTwoFAEnabled(userID int) (bool, error)
	ActivateTwoFA(userId int) error
	GenerateToken(userId int, ttl time.Duration) (string, error)
	SaveRefreshToken(UserId int, refreshToken string) error
	ValidateRefreshToken(refreshToken string) (int, error)
	BlacklistToken(token string) error
	IsTokenBlacklisted(token string) (bool, error)
}

type SendPassword interface {
	GetTokenResetPassword(email string) (int, time.Time, error)
	SaveResetToken(userID int, token string, expiry time.Time) error
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
	SendPassword  SendPassword
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
		SendPassword:  NewResetPostgres(db),
		Chat:          NewChatRepository(db.DB),
		Message:       NewMessageRepository(db),
		Notification:  NewNotificationRepository(db.DB),
		Payment:       NewPaymentRepository(db.DB),
		Subscription:  NewSubscriptionRepository(db.DB),
	}
}
