package repository

import (
	"database/sql"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) SendNotification(userID int, message string) error {
	return nil
}

func (r *NotificationRepository) GetNotifications(userID string) ([]interface{}, error) {
	return []interface{}{}, nil
}
