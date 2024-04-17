package service

import (
	"context"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) SignUp(cxt context.Context, user domain.User) error {
	return us.repo.Create(cxt, user)
}
