package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode      string
	SendTooManyError   = errors.New("发送太频繁")
	VerifyTooManyError = errors.New("验证码已失效")
)

type Codecache struct {
	cmd redis.Cmdable
}

func newCodecache(cmd redis.Cmdable) *Codecache {
	return &Codecache{
		cmd: cmd,
	}
}

/*
*

			因为直接去用业务代码查询redis中是否包含key，来判断验证码是否存在，存在并发问题
			这种并发问题，属于业务层面上，或者分布式环境下的并发问题，不是语言层面的并发问题
			并发问题场景：
		    攻击者，如果使用很多个线程，同时使用同一个号码，来触发验证码的发送，多个线程都读到
		    redis没有验证码，从而导致多个线程都会发送信息，从而造成短信成本损失
		    解决方案：
		     1 在不考虑性能的情况下：如果把key作为锁，把"查询redis执行的Get"和“发送验证码
	         之前的执行的SET”作为一个原子性的操作
	         2 就是这里利用redis单线程的机制，使用lua脚本，GET校验验证码和存储验证码,redis
	         会挨个执行每个lua脚本
*/
func (c *Codecache) Set(cxt context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(cxt, luaSetCode, []string{c.key(biz, phone), code}).Int()
	if err != nil {
		return err //调用redis出了问题
	}
	switch res {
	case -1:
		return errors.New("验证码存在，但是没有过期时间")
	case -2:
		return SendTooManyError
	default:
		return nil
	}
}

func (c *Codecache) Verify(cxt context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(cxt, luaVerifyCode, []string{c.key(biz, phone), code}).Int()
	if err != nil {
		return false, err //调用redis出了问题
	}
	switch res {
	case -1:
		return false, VerifyTooManyError
	case -2:
		return false, nil
	default:
		return true, nil
	}
}

func (c *Codecache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
