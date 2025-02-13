package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService struct {
	client *twilio.RestClient
	from   string
}

func NewSMSService(accountSID, authToken, from string) *SMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	return &SMSService{client: client, from: from}
}

func (s *SMSService) SendSMS(to, message string) error {
	if to == "" {
		return fmt.Errorf("phone number is required")
	}
	if s.from == "" {
		return fmt.Errorf("sender number is required")
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.from)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	return err
}

func (s *SMSService) GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
