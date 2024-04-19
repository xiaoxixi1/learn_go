package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository"
)

// 定义别名，进行层级传递
var (
	EmailDuplicateError   = repository.EmailDuplicateError
	InvalidPasswordOrUser = errors.New("账号或者密码不正确")
	UserNotFoundError     = repository.UserNotFoundErr
)

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

func (us *UserService) Login(cxt context.Context, email string, password string) (domain.User, error) {
	user, err := us.repo.FindByEmail(cxt, email)
	if err == repository.UserNotFoundErr {
		return domain.User{}, InvalidPasswordOrUser
	}
	if err != nil {
		return domain.User{}, err
	}
	// 验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, InvalidPasswordOrUser
	}
	return user, err
}

func (us *UserService) Edit(cxt *gin.Context, user domain.User) error {
	return us.repo.UpdateNoSensitiveInfo(cxt, user)
}

func (us *UserService) Profile(cxt *gin.Context, userid int64) (domain.User, error) {
	return us.repo.FindById(cxt, userid)
}
