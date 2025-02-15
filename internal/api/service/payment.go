package service

import (
	"github.com/polyk005/message/internal/api/repository"
)

type PaymentService struct {
	repo *repository.PaymentRepository
}

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(userID int, amount float64, purpose string) (int, error) {
	return s.repo.CreatePayment(userID, amount, purpose)
}

func (s *PaymentService) GetPaymentStatus(paymentID string) (string, error) {
	return s.repo.GetPaymentStatus(paymentID)
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(userID int, plan string) (int, error) {
	return s.repo.CreateSubscription(userID, plan)
}

func (s *SubscriptionService) GetSubscription(userID int) (map[string]interface{}, error) {
	return s.repo.GetSubscription(userID)
}

func (s *SubscriptionService) CancelSubscription(subscriptionID int) error {
	return s.repo.CancelSubscription(subscriptionID)
}
