package freecache

import (
	"context"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"log"
	"project_go/webbook/internal/repository/cache"
	"strconv"
	"sync"
	"time"
)

type CodeFreeCache struct {
	cache *freecache.Cache
	lock  sync.Mutex
}

func (c *CodeFreeCache) Set(cxt context.Context, biz, phone, code string) error {
	/**
	  获取key：
	    key存在：但是不存在过期时间 ->报错,我理解本地缓存，没有客户端是不是不存在这种场景
	    key不存在或者key发送的时间超过了1分钟 -> 可重新发送验证码
	    key存在，但是发送的时间没有超过1分钟 -> 报错发送过于频繁
	*/
	key := cache.Key(biz, phone)
	c.lock.Lock()
	defer c.lock.Unlock()
	_, expireAt, err := c.cache.GetWithExpiration([]byte(key))
	log.Println("expireAt:", expireAt)
	isTimeOut := time.Now().Add(time.Minute * 9).After(time.Unix(int64(expireAt), 0))
	if err != nil && !errors.Is(err, freecache.ErrNotFound) {
		return err // 系统错误
	}
	if errors.Is(err, freecache.ErrNotFound) || isTimeOut {
		// 发送验证码，设置有效时间
		cntkey := c.cntKey(key)
		err = c.cache.Set([]byte(key), []byte(code), 600)
		if err != nil {
			return err
		}
		return c.cache.Set([]byte(cntkey), []byte{'3'}, 600)

	}
	return cache.SendTooManyError // key存在但是发送时间没有超过1分钟
}

func (c *CodeFreeCache) Verify(cxt context.Context, biz, phone, code string) (bool, error) {
	key := cache.Key(biz, phone)
	cntket := c.cntKey(key)
	c.lock.Lock()
	defer c.lock.Unlock()
	cnt, expireAt, err := c.cache.GetWithExpiration([]byte(cntket))
	if err != nil {
		return false, err
	}
	cntValue := int(cnt[0])
	if cntValue <= 0 {
		return false, cache.VerifyTooManyError //验证次数已耗尽
	}
	value, err := c.cache.Get([]byte(key))
	if err != nil {
		return false, err
	}
	if string(value) == code {
		return true, nil // 验证正确
	}
	newCnt := strconv.Itoa(cntValue - 1)
	newExpire := time.Unix(int64(expireAt), 0).Sub(time.Now()).Microseconds()
	c.cache.Set([]byte(cntket), []byte(newCnt), int(newExpire))
	// 否则，验证码不正确
	return false, nil

}

func (c *CodeFreeCache) cntKey(key string) string {
	return fmt.Sprintf("%s:cnt", key)
}

func NewCodeFreeCache(cache *freecache.Cache) cache.CodeCache {
	return &CodeFreeCache{
		cache: cache,
	}
}
