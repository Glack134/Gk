package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/message/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (country, username, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id`, usersTable)
	row := r.db.QueryRow(query, user.Сountry, user.Username, user.Password, user.Email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) SaveResetToken(userID int, token string, expiry time.Time) error {
	query := "INSERT INTO reset_tokens (user_id, token, expiry) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, userID, token, expiry)
	return err
}

func (r *AuthPostgres) GetUser(email, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	return user, err
}

func (r *AuthPostgres) GetUserByPhone(phone string) (string, error) {
	var email string
	query := fmt.Sprintf("SELECT email FROM %s WHERE phone = $1", usersTable)
	err := r.db.Get(&email, query, phone)
	return email, err
}

func (r *AuthPostgres) GetUserByEmail(email string) (string, error) {
	var phone string
	query := fmt.Sprintf("SELECT phone FROM %s WHERE email = $1", usersTable)
	err := r.db.Get(&phone, query, email)
	return phone, err
}

func (r *AuthPostgres) GetTokenResetPassword(email string) (int, time.Time, error) {
	var userID int
	var lastSent sql.NullTime
	query := fmt.Sprintf("SELECT id, last_sent FROM %s WHERE email = $1", usersTable)
	err := r.db.QueryRow(query, email).Scan(&userID, &lastSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, fmt.Errorf("пользователь с таким email не найден")
		}
		return 0, time.Time{}, err
	}

	if lastSent.Valid {
		return userID, lastSent.Time, nil
	}

	return userID, time.Time{}, nil
}

// Получаем userID по токену
func (r *AuthPostgres) GetUserIDByToken(token string) (int, error) {
	var userID int
	query := "SELECT user_id FROM reset_tokens WHERE token=$1 AND expiry > NOW()"
	err := r.db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Обновляем пароль пользователя по userID
func (r *AuthPostgres) UpdatePasswordUserByID(userID int, newPasswordHash string) error {
	query := "UPDATE users SET password_hash=$1 WHERE id=$2"
	_, err := r.db.Exec(query, newPasswordHash, userID)
	return err
}

func (r *AuthPostgres) UpdatePasswordUser(username, newPasswordHash string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("UPDATE %s SET password_hash = $1 WHERE username = $2 RETURNING id, username, email", usersTable)
	err := r.db.QueryRow(query, newPasswordHash, username).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *AuthPostgres) MarkTokenAsUsed(token string) error {
	query := "UPDATE reset_tokens SET used = TRUE WHERE token=$1"
	_, err := r.db.Exec(query, token)
	return err
}

func (r *AuthPostgres) IsTokenUsed(token string) (bool, error) {
	var isUsed bool
	query := "SELECT used FROM reset_tokens WHERE token = $1"
	err := r.db.QueryRow(query, token).Scan(&isUsed)
	if err != nil {
		return false, err
	}
	return isUsed, nil
}

func (r *AuthPostgres) GetLastSentTime(token string) (time.Time, error) {
	var lastSentAt time.Time
	err := r.db.QueryRow("SELECT last_sent_at FROM reset_tokens WHERE token = $1", token).Scan(&lastSentAt)
	if err != nil {
		return time.Time{}, err
	}
	return lastSentAt, nil
}

//2fa

func (r *AuthPostgres) UpdateTwoFASecret(userID int, secret string) error {
	query := "UPDATE users SET two_fa_secret = $1, is_two_fa_enabled = TRUE WHERE id = $2"
	_, err := r.db.Exec(query, secret, userID)
	return err
}

func (r *AuthPostgres) GetTwoFASecret(userID int) (string, error) {
	var secret string
	query := "SELECT two_fa_secret FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userID).Scan(&secret)
	return secret, err
}
