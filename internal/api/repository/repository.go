package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
	CreateUser(phone, email string) error
	GetUserByPhone(phone string) (string, error)
	GetUserByEmail(email string) (string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
