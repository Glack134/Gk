package service

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailService(from, password, smtpHost, smtpPort string) *EmailService {
	return &EmailService{from: from, password: password, smtpHost: smtpHost, smtpPort: smtpPort}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
	return smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, []byte(msg))
}
