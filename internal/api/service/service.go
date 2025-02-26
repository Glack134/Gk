package service

import (
	"github.com/polyk005/message/internal/api/repository"
	"github.com/polyk005/message/internal/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	HashPassword(password string) (string, error)
	UpdatePasswordUserToken(token, newPassword string) error
	CheckToken(token string) error
	EnableTwoFA(userID int) (string, error)
	VerifyTwoFACode(userID int, code string) (bool, error)
	DisableTwoFA(userId int) error
	ConfirmTwoFA(userId int, code string) error
	IsTwoFAEnabled(userID int) (bool, error)
	GetUser(email, password string, checkPassword bool) (model.User, error)
}

type SendPassword interface {
	CreateResetToken(email string) (string, error)
	sendEmail(from string, password string, to string, subject string, body string) error
}

type User interface {
	GetUserProfile(userID int) (*model.User, error)
	UpdateUserProfile(user *model.User_update) error
	ValidateResetCode(code string) (string, error)
	UpdateUserEmail(userID int, newEmail string) error
}

type Chat interface {
	CreateChat(chatName string, userIDs ...int) (int, error)
	AddParticipant(chatID, userID int) error
	GetChatsForUser(userID int) ([]map[string]interface{}, error)
	UserExists(username string) (bool, error)
	ChatExistsBetweenUsers(userID1, userID2 int) (int, error)
	GetUserIDByUsername(username string) (int, error)
	DeleteChatForAll(chatID int) error
	DeleteChatForUser(chatID, userID int) error
	FindExistingChat(userIDs []int) (int, error)
}

type Message interface {
	GetMessages(chatID string) ([]repository.Message, error)
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
	Authorization
	SendPassword
	User
	Chat
	Message
	Notification
	Payment
	Subscription
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		SendPassword:  NewSendPassword(repos.SendPassword),
		User:          NewUserService(repos.User),
		Chat:          NewChatService(repos.Chat),
		Message:       NewMessageService(repos.Message),
		Notification:  NewNotificationService(repos.Notification),
		Payment:       NewPaymentService(repos.Payment),
		Subscription:  NewSubscriptionService(repos.Subscription),
	}
}
