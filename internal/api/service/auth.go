package service

import (
	"fmt"

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
