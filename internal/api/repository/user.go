package repository

import (
	"database/sql"

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

func (r *UserRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET username = $1 WHERE id = $1`
	_, err := r.db.Exec(query, user.Username, user.Id)
	return err
}
