package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(phone, email string) error {
	query := fmt.Sprintf("INSERT INTO %s (phone, email) VALUES ($1, $2)", usersTable)
	_, err := r.db.Exec(query, phone, email)
	return err
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
