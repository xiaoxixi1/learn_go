package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository"
	repomocks "project_go/webbook/internal/repository/mocks"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#hello")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#hello"))
	assert.NoError(t, err)
}

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 预期输入
		cxt      context.Context
		email    string
		password string

		// 预期输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$kVlNjLYYwV6j0A2krcLHb.hsrnnskbUA6RPqInTf.I2mClfxXO8Xm",
					Phone:    "1253553",
				}, nil)
				return userRepo
			},
			cxt:      nil,
			email:    "123@qq.com",
			password: "123456#hello",
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$kVlNjLYYwV6j0A2krcLHb.hsrnnskbUA6RPqInTf.I2mClfxXO8Xm",
				Phone:    "1253553",
			},
		},
		{
			name: "用户名不正确",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.UserNotFoundErr)
				return userRepo
			},
			cxt:      nil,
			email:    "123@qq.com",
			password: "123456#hello",
			wantUser: domain.User{},
			wantErr:  InvalidPasswordOrUser,
		},
		{
			name: "密码不正确",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$balalalsls",
					Phone:    "1253553",
				}, nil)
				return userRepo
			},
			cxt:      nil,
			email:    "123@qq.com",
			password: "123456#hello",
			wantUser: domain.User{},
			wantErr:  InvalidPasswordOrUser,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("db error"))
				return userRepo
			},
			cxt:      nil,
			email:    "123@qq.com",
			password: "123456#hello",
			wantUser: domain.User{},
			wantErr:  errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRepo := tc.mock(ctrl)
			userSvc := NewUserService(userRepo)
			user, err := userSvc.Login(tc.cxt, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
