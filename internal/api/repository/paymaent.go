package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

type PaymentDetails struct {
	Amount   float64 `db:"amount"`
	Currency string  `db:"currency"`
}

func (r *PaymentRepository) GetPaymentDetails(paymentID int) (*PaymentDetails, error) {
	var details PaymentDetails
	query := `SELECT amount, currency FROM payments WHERE id = $1`
	err := r.db.QueryRow(query, paymentID).Scan(&details.Amount, &details.Currency)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

func (r *PaymentRepository) CreatePayment(userID int, amount float64, purpose, paymentMethod string) (int, error) {
	var paymentID int
	query := `INSERT INTO payments (user_id, amount, purpose, payment_method, status) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, userID, amount, purpose, paymentMethod, "pending").Scan(&paymentID)
	return paymentID, err
}

func (r *PaymentRepository) GetPaymentStatus(ctx context.Context, paymentID string) (string, error) {
	var status string
	query := `SELECT status FROM payments WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(paymentID int, status string) error {
	query := `UPDATE payments SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, paymentID)
	return err
}

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) CreateSubscription(userID int, plan string) (int, error) {
	var subscriptionID int
	query := `INSERT INTO subscriptions (user_id, plan, status, start_date, end_date) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	startDate := time.Now()
	endDate := startDate.AddDate(0, 1, 0) // Подписка на 1 месяц
	err := r.db.QueryRow(query, userID, plan, "active", startDate, endDate).Scan(&subscriptionID)
	return subscriptionID, err
}

func (r *SubscriptionRepository) GetSubscription(userID int) (map[string]interface{}, error) {
	var (
		id        int
		plan      string
		status    string
		startDate time.Time
		endDate   time.Time
	)

	query := `SELECT id, plan, status, start_date, end_date FROM subscriptions WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&id, &plan, &status, &startDate, &endDate)
	if err != nil {
		return nil, err
	}

	subscription := map[string]interface{}{
		"id":         id,
		"plan":       plan,
		"status":     status,
		"start_date": startDate,
		"end_date":   endDate,
	}

	return subscription, nil
}

func (r *SubscriptionRepository) CancelSubscription(subscriptionID int) error {
	query := `UPDATE subscriptions SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, "cancelled", subscriptionID)
	return err
}
