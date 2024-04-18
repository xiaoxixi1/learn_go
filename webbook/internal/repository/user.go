package repository

import (
	"context"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository/dao"
)

var EmailDuplicateError = dao.EmailDuplicateError

type UserRepository struct {
	dao *dao.UserDao
}

func NewUseRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}
func (repo *UserRepository) Create(cxt context.Context, user domain.User) error {
	return repo.dao.Insert(cxt, &dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}