package service

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/polyk005/message/internal/api/repository"
)

type AuthService struct {
	repo      repository.Authorization
	sms       *SMSService
	codeStore map[string]string
}

func NewAuthService(repo repository.Authorization, sms *SMSService) *AuthService {
	return &AuthService{repo: repo, sms: sms, codeStore: make(map[string]string)}
}

func (s *AuthService) SendVerificationCode(identifier string) error {
	code := s.sms.GenerateCode()
	s.codeStore[identifier] = code

	if isPhone(identifier) {
		return s.sms.SendSMS(identifier, "Your verification code is: "+code)
	} else if isEmail(identifier) {
		return fmt.Errorf("email sending not implemented yet")
	}
	return fmt.Errorf("invalid identifier")
}

func (s *AuthService) VerifyCode(identifier, code string) bool {
	storedCode, ok := s.codeStore[identifier]
	if !ok {
		return false
	}
	return storedCode == code
}

func (s *AuthService) SignUp(identifier, code string) error {
	if !s.VerifyCode(identifier, code) {
		return fmt.Errorf("invalid verification code")
	}

	if isPhone(identifier) {
		return s.repo.CreateUser(identifier, "")
	} else if isEmail(identifier) {
		return s.repo.CreateUser("", identifier)
	}
	return fmt.Errorf("invalid identifier")
}

func (s *AuthService) SignIn(identifier, code string) error {
	if !s.VerifyCode(identifier, code) {
		return fmt.Errorf("invalid verification code")
	}
	return nil
}

func isPhone(identifier string) bool {
	return len(identifier) >= 10
}

func isEmail(identifier string) bool {
	return len(identifier) >= 5 && identifier[len(identifier)-4:] == ".com"
}

func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your-secret-key"), nil // Замените на ваш секретный ключ
	})
	if err != nil {
		return 0, err
	}

	// Извлекаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64)) // Извлекаем userID из токена
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
