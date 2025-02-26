package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/polyk005/message/internal/api/repository"
	"github.com/polyk005/message/internal/model"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt        = "6765fgvbhhgf35vfu9jft5tg"
	signingKey  = "qjvkvnsjdnj2njn29njv**@9un19@!33"
	resetingKey = "fa#dh$bsia1*&2rffvsv2135v#eg*#"
	tokenTTL    = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func (s *AuthService) UpdatePasswordUserToken(token, newPassword string) error {
	isUsed, err := s.repo.IsTokenUsed(token)
	if err != nil {
		return err
	}
	if isUsed {
		return errors.New("token has already been used")
	}

	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return err
	}

	newPasswordHash := s.generatePasswordHash(newPassword)

	err = s.repo.UpdatePasswordUserByID(userID, newPasswordHash)
	if err != nil {
		return err
	}

	err = s.repo.MarkTokenAsUsed(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckToken(token string) error {
	isUsed, err := s.repo.IsTokenUsed(token)
	if err != nil {
		return err
	}
	if isUsed {
		return errors.New("token has already been used")
	}

	return nil
}

func (s *AuthService) GetUser(email, password string) (model.User, error) {
	// Вызываем метод GetUser  из репозитория
	user, err := s.repo.GetUser(email, password)
	if err != nil {
		return model.User{}, err // Возвращаем ошибку, если произошла ошибка при получении пользователя
	}

	// Проверяем, существует ли пользователь и совпадают ли email и пароль
	if user.Id == 0 {
		return model.User{}, fmt.Errorf("invalid email or password") // Возвращаем ошибку, если пользователь не найден
	}

	return user, nil // Возвращаем пользователя, если всё в порядке
}

// 2fa
func (s *AuthService) EnableTwoFA(userID int) (string, error) {
	// Генерируем новый секрет TOTP
	secret, err := totp.Generate(totp.GenerateOpts{
		AccountName: "2FA", // Название вашего приложения
		Issuer:      "Gk",  // Название вашего издателя
	})
	if err != nil {
		return "", err
	}

	// Сохраняем секрет в базе данных
	err = s.repo.UpdateTwoFASecret(userID, secret.Secret())
	if err != nil {
		return "", err
	}

	return secret.URL(), nil
}

func (s *AuthService) VerifyTwoFACode(userID int, code string) (bool, error) {
	secret, err := s.repo.GetTwoFASecret(userID)
	if err != nil {
		return false, err
	}

	return totp.Validate(code, secret), nil
}

func (s *AuthService) DisableTwoFA(userId int) error {
	isEnabled, err := s.repo.IsTwoFAEnabled(userId)
	if err != nil {
		return err
	}

	if !isEnabled {
		return errors.New("Two-Factor Authentication is not enabled")
	}

	return s.repo.DisableTwoFA(userId)
}

func (s *AuthService) ConfirmTwoFA(userId int, code string) error {
	secret, err := s.repo.GetTwoFASecret(userId)
	if err != nil {
		return fmt.Errorf("failed to get 2FA secret: %w", err)
	}

	isValid := validateTwoFACode(secret, code)
	if !isValid {
		return fmt.Errorf("invalid confirmation code")
	}

	err = s.repo.ActivateTwoFA(userId)
	if err != nil {
		return fmt.Errorf("failed to activate 2FA: %w", err)
	}

	return nil
}

func validateTwoFACode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (s *AuthService) IsTwoFAEnabled(userID int) (bool, error) {
	return s.repo.IsTwoFAEnabled(userID)
}
