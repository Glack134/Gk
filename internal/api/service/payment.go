package service

import (
	"context"
	"strings"

	"github.com/polyk005/message/internal/api/repository"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

type PaymentDetails struct {
	Amount   float64
	Currency string
}

type PaymentService struct {
	repo *repository.PaymentRepository
}

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(userID int, amount float64, purpose, paymentMethod, currency string) (int, error) {
	// Создаем платеж в Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)), // Сумма в центах
		Currency: stripe.String(strings.ToLower(currency)),
	}
	_, err := paymentintent.New(params) // Используем _, если pi не нужна
	if err != nil {
		return 0, err
	}

	// Сохраняем платеж в базе данных
	paymentID, err := s.repo.CreatePayment(userID, amount, purpose, paymentMethod, currency)
	if err != nil {
		return 0, err
	}

	return paymentID, nil
}

func (s *PaymentService) GetPaymentStatus(ctx context.Context, paymentID string) (string, error) {
	return s.repo.GetPaymentStatus(ctx, paymentID)
}

func (s *PaymentService) GetPaymentID(userID int, amount float64, purpose string) (int, error) {
	return s.repo.GetPaymentID(userID, amount, purpose)
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(userID int, plan string, paymentID int) (int, error) {
	return s.repo.CreateSubscription(userID, plan, paymentID)
}

func (s *SubscriptionService) GetSubscription(userID int) (map[string]interface{}, error) {
	return s.repo.GetSubscription(userID)
}

func (s *PaymentService) GetPaymentDetails(paymentID int) (*repository.PaymentDetails, error) {
	return s.repo.GetPaymentDetails(paymentID)
}

func (s *PaymentService) UpdatePaymentStatus(paymentID int, status string) error {
	return s.repo.UpdatePaymentStatus(paymentID, status)
}

func (s *SubscriptionService) CancelSubscription(subscriptionID int) error {
	return s.repo.CancelSubscription(subscriptionID)
}
