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

type UserRepository interface {
	Create(cxt context.Context, user domain.User) error
	FindByEmail(cxt context.Context, email string) (domain.User, error)
	UpdateNoSensitiveInfo(cxt context.Context, user domain.User) error
	FindById(cxt context.Context, userid int64) (domain.User, error)
	FindByPhone(cxt context.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	dao dao.UserDao
	c   cache.UserCache
}

// NewUserRepositoryV2 强耦合到了 JSON
//func NewUserRepositoryV2(cfgBytes string) *CachedUserRepository {
//	var cfg DBConfig
//	err := json.Unmarshal([]byte(cfgBytes), &cfg)
//}

// NewUserRepositoryV1 强耦合（跨层的），严重缺乏扩展性
//func NewUserRepositoryV1(dbCfg DBConfig, cCfg CacheConfig) (*CachedUserRepository, error) {
//	db, err := gorm.Open(mysql.Open(dbCfg.DSN))
//	if err != nil {
//		return nil, err
//	}
//	ud := dao.NewUserDAO(db)
//	uc := cache.NewUserCache(redis.NewClient(&redis.Options{
//		Addr: cCfg.Addr,
//	}))
//	return &CachedUserRepository{
//		dao:   ud,
//		cache: uc,
//	}, nil
//}

// 依赖注入的写法
func NewUseRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao: dao,
		c:   cache,
	}
}
func (repo *CachedUserRepository) Create(cxt context.Context, user domain.User) error {
	return repo.dao.Insert(cxt, repo.toEntity(user))
}

func (repo *CachedUserRepository) FindByEmail(cxt context.Context, email string) (domain.User, error) {
	us, err := repo.dao.QueryByEmail(cxt, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(us), nil
}

func (repo *CachedUserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:              user.Id,
		Email:           user.Email.String,
		Phone:           user.Phone.String,
		Name:            user.Name,
		Password:        user.Password,
		Birthday:        time.UnixMilli(user.Birthday),
		PersonalProfile: user.PersonalProfile,
		CTime:           time.UnixMilli(user.CTime),
		UTime:           time.UnixMilli(user.UTime),
	}
}

func (repo *CachedUserRepository) toEntity(user domain.User) *dao.User {
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
		CTime:           user.CTime.UnixMilli(),
		UTime:           user.UTime.UnixMilli(),
	}
}

func (repo *CachedUserRepository) UpdateNoSensitiveInfo(cxt context.Context, user domain.User) error {
	return repo.dao.Update(cxt, dao.User{
		Id:              user.Id,
		Name:            user.Name,
		Birthday:        user.Birthday.UnixMilli(),
		PersonalProfile: user.PersonalProfile,
	})
}

func (repo *CachedUserRepository) FindById(cxt context.Context, userid int64) (domain.User, error) {
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

func (repo *CachedUserRepository) FindByPhone(cxt context.Context, phone string) (domain.User, error) {
	us, err := repo.dao.QueryByPhone(cxt, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(us), err
}
