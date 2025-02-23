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

func (s *UserService) GetUserProfile(userID int) (*model.User, error) {
	return s.repo.GetUserID(userID)
}

func (s *UserService) UpdateUserProfile(user *model.User_update) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) UpdateUserEmail(userID int, newEmail string) error {
	return s.repo.UpdateUserEmail(userID, newEmail)
}

func (s *UserService) ValidateResetCode(code string) (string, error) {
	return s.ValidateResetCode(code)
}
