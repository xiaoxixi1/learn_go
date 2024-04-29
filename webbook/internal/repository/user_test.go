package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository/cache"
	cachemocks "project_go/webbook/internal/repository/cache/mocks"
	"project_go/webbook/internal/repository/dao"
	daomocks "project_go/webbook/internal/repository/dao/mocks"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	nowMili := time.Now().UnixMilli()
	now := time.UnixMilli(nowMili)
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)

		cxt    context.Context
		userId int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(nil, int64(1)).Return(domain.User{}, errors.New("缓存未命中"))
				userDao := daomocks.NewMockUserDao(ctrl)
				userDao.EXPECT().QueryById(nil, int64(1)).Return(dao.User{
					Id: int64(1),
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Phone: sql.NullString{
						String: "1233",
						Valid:  true,
					},
					Name:            "xixi",
					Birthday:        int64(nowMili),
					PersonalProfile: "自我介绍",
					Password:        "123",
					// 创建时间，使用UTC 0的毫秒数，时区的转换一般统一让前端转换，或者留到要传给前端时转换
					CTime: nowMili,
					UTime: nowMili,
				}, nil)
				userCache.EXPECT().Set(nil, domain.User{
					Id:              int64(1),
					Email:           "123@qq.com",
					Password:        "123",
					Name:            "xixi", //昵称
					Birthday:        now,
					PersonalProfile: "自我介绍", // 个人简介
					Phone:           "1233",
					CTime:           now,
					UTime:           now,
				})
				return userCache, userDao
			},
			cxt:    nil,
			userId: int64(1),
			wantUser: domain.User{
				Id:              int64(1),
				Email:           "123@qq.com",
				Password:        "123",
				Name:            "xixi", //昵称
				Birthday:        now,
				PersonalProfile: "自我介绍", // 个人简介
				Phone:           "1233",
				CTime:           now,
				UTime:           now,
			},
		},
		{
			name: "查找成功，缓存命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(nil, int64(1)).Return(domain.User{
					Id:       int64(1),
					Password: "123",
				}, nil)
				userDao := daomocks.NewMockUserDao(ctrl)
				//userDao.EXPECT().QueryById(nil,1,).Return()
				return userCache, userDao
			},
			cxt:    nil,
			userId: int64(1),
			wantUser: domain.User{
				Id:       int64(1),
				Password: "123",
			},
		},
		{
			name: "数据库查找失败，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(nil, int64(1)).Return(domain.User{}, errors.New("缓存未命中"))
				userDao := daomocks.NewMockUserDao(ctrl)
				userDao.EXPECT().QueryById(nil, int64(1)).Return(dao.User{}, errors.New("db err"))
				return userCache, userDao
			},
			cxt:      nil,
			userId:   int64(1),
			wantUser: domain.User{},
			wantErr:  errors.New("db err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userCache, userDao := tc.mock(ctrl)
			userRepo := NewUseRepository(userDao, userCache)
			user, err := userRepo.FindById(tc.cxt, tc.userId)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
