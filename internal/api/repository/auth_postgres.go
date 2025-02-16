package repository

import (
	"fmt"

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

	query := fmt.Sprintf(`INSERT INTO %s (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, user.Username, user.Password, user.Email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
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
