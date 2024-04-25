package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository/cache"
	"project_go/webbook/internal/repository/dao"
	"time"
)

var (
	EmailDuplicateError = dao.EmailDuplicateError
	UserNotFoundErr     = dao.RecordNotFoundErr
)

type UserRepository struct {
	dao *dao.UserDao
	c   *cache.UserCache
}

func NewUseRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao: dao,
		c:   cache,
	}
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
		Birthday:        time.UnixMilli(user.Birthday),
		PersonalProfile: user.PersonalProfile,
	}
}

func (repo *UserRepository) UpdateNoSensitiveInfo(cxt *gin.Context, user domain.User) error {
	return repo.dao.Update(cxt, dao.User{
		Id:              user.Id,
		Name:            user.Name,
		Birthday:        user.Birthday.UnixMilli(),
		PersonalProfile: user.PersonalProfile,
	})
}

func (repo *UserRepository) FindById(cxt *gin.Context, userid int64) (domain.User, error) {
	du, err := repo.c.Get(cxt, userid)
	// 缓存命中
	if err == nil {
		return du, err
	}
	us, err := repo.dao.QueryById(cxt, userid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(us)
	err = repo.c.Set(cxt, du)
	return du, nil
}
