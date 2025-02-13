package service

import (
	"fmt"

	"github.com/polyk005/message/internal/api/repository"
)

type AuthService struct {
	repo      repository.Authorization
	sms       *SMSService
	codeStore map[string]string // временное хранилище кодов
}

func NewAuthService(repo repository.Authorization, sms *SMSService) *AuthService {
	return &AuthService{repo: repo, sms: sms, codeStore: make(map[string]string)}
}

func (s *AuthService) SendVerificationCode(phone string) error {
	code := s.sms.GenerateCode()
	s.codeStore[phone] = code
	return s.sms.SendSMS(phone, "Your verification code is: "+code)
}

func (s *AuthService) VerifyCode(phone, code string) bool {
	storedCode, ok := s.codeStore[phone]
	if !ok {
		return false
	}
	return storedCode == code
}

func (s *AuthService) SignUp(phone, code string) error {
	if !s.VerifyCode(phone, code) {
		return fmt.Errorf("invalid verification code")
	}
	// Логика создания пользователя в базе данных
	return nil
}

func (s *AuthService) SignIn(phone, code string) error {
	if !s.VerifyCode(phone, code) {
		return fmt.Errorf("invalid verification code")
	}
	// Логика авторизации пользователя
	return nil
}
