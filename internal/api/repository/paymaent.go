package repository

import (
	"database/sql"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(userID int, amount float64, purpose string) (int, error) {
	var paymentID int
	query := `INSERT INTO payments (user_id, amount, purpose, status) 
	          VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, userID, amount, purpose, "pending").Scan(&paymentID)
	return paymentID, err
}

func (r *PaymentRepository) GetPaymentStatus(paymentID string) (string, error) {
	var status string
	query := `SELECT status FROM payments WHERE id = $1`
	err := r.db.QueryRow(query, paymentID).Scan(&status)
	return status, err
}

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
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
