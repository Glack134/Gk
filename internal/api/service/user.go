package service

import (
	"github.com/polyk005/message/internal/api/repository"
	"github.com/polyk005/message/internal/model"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUserProfile возвращает профиль пользователя
func (s *UserService) GetUserProfile(userID int) (*model.User, error) {
	return s.repo.GetUserID(userID)
}

// UpdateUserProfile обновляет профиль пользователя
func (s *UserService) UpdateUserProfile(user *model.User) error {
	return s.repo.UpdateUser(user)
}
