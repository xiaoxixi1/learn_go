package ratelimit

import (
	"context"
	"errors"
	"project_go/webbook/internal/service/sms"
	"project_go/webbook/pkg/limiter"
)

// 装饰器设计模式，在不改变原来短信发送实现的情况下，增加限流器的功能
type RateLimitSMSService struct {
	// 被装饰者
	svc   sms.Service
	limit limiter.Limiter
	key   string
}

var errLimited = errors.New("触发限流")

func (r RateLimitSMSService) Send(cxt context.Context, tplId string, args []string, number ...string) error {
	limited, err := r.limit.Limit(cxt, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	// 最终调用被装饰者的方法
	return r.svc.Send(cxt, tplId, args, number...)
}

func NewRateLimitSMSService(svc sms.Service, limit limiter.Limiter, key string) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:   svc,
		limit: limit,
		key:   key,
	}
}
