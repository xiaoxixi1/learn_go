package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository"
)

// 定义别名，进行层级传递
var (
	UserDuplicateError    = repository.UserDuplicateError
	InvalidPasswordOrUser = errors.New("账号或者密码不正确")
	UserNotFoundError     = repository.UserNotFoundErr
)

/*
*

	面向接口编程是指将应用程序定义为组件的集合，组件和组件之间的通信必须通过接口
	如果要用到另外一个类型，那么你肯定用的是接口
	结合依赖注入，如果A调用B，B是一个接口，而后在初始化的时候，注入了一个实现了B接口的实例
*/
type UserService interface {
	SignUp(cxt context.Context, user domain.User) error
	Login(cxt context.Context, email string, password string) (domain.User, error)
	Edit(cxt context.Context, user domain.User) error
	Profile(cxt context.Context, userid int64) (domain.User, error)
	FindOrCreate(cxt context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(cxt context.Context, wechatInfo domain.WeChatDomain) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) SignUp(cxt context.Context, user domain.User) error {
	pHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(pHash)
	return us.repo.Create(cxt, user)
}

func (us *userService) Login(cxt context.Context, email string, password string) (domain.User, error) {
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
		return domain.User{}, InvalidPasswordOrUser
	}
	return user, err
}

func (us *userService) Edit(cxt context.Context, user domain.User) error {
	return us.repo.UpdateNoSensitiveInfo(cxt, user)
}

func (us *userService) Profile(cxt context.Context, userid int64) (domain.User, error) {
	return us.repo.FindById(cxt, userid)
}

func (us *userService) FindOrCreate(cxt context.Context, phone string) (domain.User, error) {
	user, err := us.repo.FindByPhone(cxt, phone)
	if err != UserNotFoundError {
		// 如果不是没找到，是nil或者是别的报错，都是直接返回
		return user, err
	}
	// 否则就是用户没有找到，需要注册
	err = us.repo.Create(cxt, domain.User{
		Phone: phone,
	})
	// 如果是别的报错，则直接返回
	if err != nil && err != UserDuplicateError {
		return domain.User{}, err
	}
	// 如果是用户已经存在，或者没有错误，则再查一遍
	return us.repo.FindByPhone(cxt, phone)
}

func (us *userService) FindOrCreateByWechat(cxt context.Context, wechatInfo domain.WeChatDomain) (domain.User, error) {
	user, err := us.repo.FindByOpenId(cxt, wechatInfo.OpenId)
	if err != UserNotFoundError {
		// 如果不是没找到，是nil或者是别的报错，都是直接返回
		return user, err
	}
	// 否则就是用户没有找到，需要注册
	err = us.repo.Create(cxt, domain.User{
		WechatInfo: wechatInfo,
	})
	// 如果是别的报错，则直接返回
	if err != nil && err != UserDuplicateError {
		return domain.User{}, err
	}
	// 如果是用户已经存在，或者没有错误，则再查一遍
	return us.repo.FindByOpenId(cxt, wechatInfo.OpenId)
}
