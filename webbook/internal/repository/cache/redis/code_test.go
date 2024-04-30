package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project_go/webbook/internal/repository/cache"
	"project_go/webbook/internal/repository/cache/redismocks"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) (cmd redis.Cmdable)
		ctx   context.Context
		biz   string
		phone string
		code  string

		wantError error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))
				cmd.SetErr(nil)
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{cache.Key("UserLogin", "123")}, []any{"123"}).Return(cmd)
				return res
			},
			ctx:       context.Background(),
			biz:       "UserLogin",
			phone:     "123",
			code:      "123",
			wantError: nil,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))
				cmd.SetErr(errors.New("redis错误"))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{cache.Key("UserLogin", "123")}, []any{"123"}).Return(cmd)
				return res
			},
			ctx:       nil,
			biz:       "UserLogin",
			phone:     "123",
			code:      "123",
			wantError: errors.New("redis错误"),
		},
		{
			name: "验证码存在，没有设置过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-1))
				cmd.SetErr(errors.New("验证码存在，但是没有过期时间"))
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{cache.Key("UserLogin", "123")}, []any{"123"}).Return(cmd)
				return res
			},
			ctx:       nil,
			biz:       "UserLogin",
			phone:     "123",
			code:      "123",
			wantError: errors.New("验证码存在，但是没有过期时间"),
		},
		{
			name: "发送过于频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-2))
				cmd.SetErr(cache.SendTooManyError)
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{cache.Key("UserLogin", "123")}, []any{"123"}).Return(cmd)
				return res
			},
			ctx:       nil,
			biz:       "UserLogin",
			phone:     "123",
			code:      "123",
			wantError: cache.SendTooManyError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cmd := tc.mock(ctrl)
			codeCache := NewCodeCache(cmd)
			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantError, err)
		})
	}
}
