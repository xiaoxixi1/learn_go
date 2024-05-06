package failover

import (
	"context"
	"project_go/webbook/internal/service/sms"
	"sync/atomic"
)

type TimeOutFailOverSMSService struct {
	svcs []sms.Service
	// 当前使用得服务商
	idx int32
	// 连续超时次数
	cnt int32
	// 切换的阈值
	threshold int32
}

func NewTimeOutFailOverSMSService(svcs []sms.Service, threshold int32) *TimeOutFailOverSMSService {
	return &TimeOutFailOverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeOutFailOverSMSService) Send(cxt context.Context, tplId string, args []string, number ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	// 超过阈值，进行切换
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		// 防止多个人进行切换，需要用原子并发工具
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	/**
	  这里的实现并不是非常严格的连续N个超时就切换，只是近似，但是不影响实际效果
	*/
	svc := t.svcs[idx]
	err := svc.Send(cxt, tplId, args, number...)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 遇到了错误，但是又不是超时错误
		// 可以增加，也可以不增加
		// 如果是EOF之类的错误，还可以考虑直接切换
	}
	return err
}
