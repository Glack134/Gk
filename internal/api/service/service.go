package service

import "github.com/polyk005/message/internal/api/repository"

type Authorization interface {
	SendVerificationCode(phone string) error
	VerifyCode(phone, code string) bool
	SignUp(phone, code string) error
	SignIn(phone, code string) error
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, sms *SMSService) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, sms),
	}
}
