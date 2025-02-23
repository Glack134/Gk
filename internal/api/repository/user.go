package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/polyk005/message/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserID(userID int) (*model.User, error) {
	var user model.User
	query := `SELECT id, country, username, email FROM users WHERE id = $1`
	err := r.db.QueryRow(query, userID).Scan(&user.Id, &user.Сountry, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User_update) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	_, err := r.db.Exec(query, user.Username, user.Id)
	return err
}

func (r *UserRepository) UpdateUserEmail(userID int, newEmail string) error {
	query := `UPDATE users SET email = $1 WHERE id = $2`
	_, err := r.db.Exec(query, newEmail, userID)
	return err
}

func (r *UserRepository) UpdateUserPasswordByEmail(email, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE email = $2`
	_, err := r.db.Exec(query, hashedPassword, email)
	return err
}

func (r *UserRepository) ValidateResetCode(code string) (string, error) {
	var email string

	// Проверка кода восстановления в базе данных
	query := `SELECT email FROM password_resets WHERE code = $1 AND expires_at > NOW()`
	err := r.db.QueryRow(query, code).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid or expired reset code")
		}
		return "", fmt.Errorf("failed to validate reset code: %w", err)
	}

	return email, nil
}

func (r *UserRepository) GetTokenResetPassword(email string) (int, time.Time, error) {
	var userID int
	var lastSent time.Time
	row := r.db.QueryRow("SELECT user_id, last_sent FROM users WHERE email = $1", email)
	err := row.Scan(&userID, &lastSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, fmt.Errorf("пользователь с таким email не найден")
		}
		return 0, time.Time{}, err
	}
	return userID, lastSent, nil
}
