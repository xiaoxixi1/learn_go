package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"project_go/webbook/internal/repository/cache"
	"testing"
	"time"
)

func TestRedisUserCache_e2e_Set(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		ctx    context.Context
		biz    string
		phone  string
		code   string

		wantError error
	}{
		{
			name: "设置成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:userLogin:12345"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, code, "123456")
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:       context.Background(),
			biz:       "userLogin",
			phone:     "12345",
			code:      "123456",
			wantError: nil,
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				//提前准备一条数据
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				// 删除数据
				_, err := rdb.Del(ctx, key).Result()
				assert.NoError(t, err)
			},
			ctx:       context.Background(),
			biz:       "userLogin",
			phone:     "12345",
			code:      "123456",
			wantError: cache.SendTooManyError,
		},
		{
			name: "没有设置过期时间",
			before: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				//提前准备一条数据
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
				defer cancle()
				key := "phone_code:userLogin:12345"
				// 删除数据
				_, err := rdb.Del(ctx, key).Result()
				assert.NoError(t, err)
			},
			ctx:       context.Background(),
			biz:       "userLogin",
			phone:     "12345",
			code:      "123456",
			wantError: errors.New("验证码存在，但是没有过期时间"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			codeCache := NewCodeCache(rdb)
			err := codeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantError, err)
		})
	}
}
