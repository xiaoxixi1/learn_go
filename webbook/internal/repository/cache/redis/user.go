package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project_go/webbook/internal/domain"
	"project_go/webbook/internal/repository/cache"
	"time"
)

type RedisUserCache struct {
	cmd        redis.Cmdable
	expireTime time.Duration
}

func NewUserCache(cmd redis.Cmdable) cache.UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expireTime: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Get(cxt context.Context, userid int64) (domain.User, error) {
	key := cache.key(userid)
	data, err := cache.cmd.Get(cxt, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	return user, err
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *RedisUserCache) Set(cxt context.Context, us domain.User) error {
	key := cache.key(us.Id)
	data, err := json.Marshal(us)
	if err != nil {
		return err
	}
	return cache.cmd.Set(cxt, key, data, cache.expireTime).Err()
}
