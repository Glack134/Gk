package service

import "github.com/polyk005/message/internal/api/repository"

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(userID int, amount float64, purpose string) (int, error) {
	return s.repo.CreatePayment(userID, amount, purpose)
}

func (s *PaymentService) GetPaymentStatus(paymentID string) (string, error) {
	return s.repo.GetPaymentStatus(paymentID)
}
