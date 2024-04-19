package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository/dao"
)

var (
	EmailDuplicateError = dao.EmailDuplicateError
	UserNotFoundErr     = dao.RecordNotFoundErr
)

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

func (repo *UserRepository) FindByEmail(cxt context.Context, email string) (domain.User, error) {
	us, err := repo.dao.QueryByEmail(cxt, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(us), nil
}

func (repo *UserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:              user.Id,
		Email:           user.Email,
		Name:            user.Name,
		Password:        user.Password,
		Birthday:        user.Birthday,
		PersonalProfile: user.PersonalProfile,
	}
}

func (repo *UserRepository) Update(cxt *gin.Context, user domain.User) error {
	return repo.dao.Update(cxt, dao.User{
		Id:              user.Id,
		Name:            user.Name,
		Birthday:        user.Birthday,
		PersonalProfile: user.PersonalProfile,
	})
}

func (repo *UserRepository) FindById(cxt *gin.Context, userid int64) (domain.User, error) {
	us, err := repo.dao.QueryById(cxt, userid)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(us), nil
}
