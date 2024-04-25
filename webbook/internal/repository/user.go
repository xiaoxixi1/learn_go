package repository

import (
	"context"
	"database/sql"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository/cache"
	"project_go/webbook/internal/repository/dao"
	"time"
)

var (
	UserDuplicateError = dao.UserDuplicateError
	UserNotFoundErr    = dao.RecordNotFoundErr
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
	return repo.dao.Insert(cxt, repo.toEntity(user))
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
		Email:           user.Email.String,
		Phone:           user.Phone.String,
		Name:            user.Name,
		Password:        user.Password,
		Birthday:        time.UnixMilli(user.Birthday),
		PersonalProfile: user.PersonalProfile,
	}
}

func (repo *UserRepository) toEntity(user domain.User) *dao.User {
	return &dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Name:            user.Name,
		Password:        user.Password,
		Birthday:        user.Birthday.UnixMilli(),
		PersonalProfile: user.PersonalProfile,
	}
}

func (repo *UserRepository) UpdateNoSensitiveInfo(cxt context.Context, user domain.User) error {
	return repo.dao.Update(cxt, dao.User{
		Id:              user.Id,
		Name:            user.Name,
		Birthday:        user.Birthday.UnixMilli(),
		PersonalProfile: user.PersonalProfile,
	})
}

func (repo *UserRepository) FindById(cxt context.Context, userid int64) (domain.User, error) {
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

func (repo *UserRepository) FindByPhone(cxt context.Context, phone string) (domain.User, error) {
	us, err := repo.dao.QueryByPhone(cxt, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(us), err
}
