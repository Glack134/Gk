package service

import "github.com/polyk005/message/internal/api/repository"

type Authorization interface {
	SendVerificationCode(identifier string) error
	VerifyCode(identifier, code string) bool
	SignUp(identifier, code string) error
	SignIn(identifier, code string) error
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, sms *SMSService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, sms),
	}
}
