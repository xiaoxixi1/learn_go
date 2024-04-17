package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository"
)

// 定义别名，进行层级传递
var EmailDuplicateError = repository.EmailDuplicateError

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) SignUp(cxt context.Context, user domain.User) error {
	pHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(pHash)
	return us.repo.Create(cxt, user)
}
