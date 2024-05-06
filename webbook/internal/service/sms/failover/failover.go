package failover

import (
	"context"
	"errors"
	"log"
	"project_go/webbook/internal/service/sms"
	"sync/atomic"
)

type FailOverSMSService struct {
	svcs []sms.Service
	/**
	  V1版本进行优化
	    1 其实svc 动态计算
	    2 区别了错误：context.DeadlineExceeded 和context.Cancled，都是跟用户密切相馆的，所以直接返回
	*/
	idx uint64
}

/*
*

	轮询所有的可用服务商，全部轮询完了都没成功，说明所有的服务器都挂了，更多可能是本地网络挂了
	缺点：
	   每次从头开始轮询，绝大多数请求会在svc[0]就成功，负载不均衡
	   如果svcs有几十个，轮询都很慢
*/
func (f FailOverSMSService) Send(cxt context.Context, tplId string, args []string, number ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(cxt, tplId, args, number...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("尝试了所有的服务商，但是都发送失败了")
}

func (f FailOverSMSService) SendV1(cxt context.Context, tplId string, args []string, number ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(cxt, tplId, args, number...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			// 调用者设置的超市时间到了，或者调用者直接取消了，就直接返回
			return err
		}
		// 其他情况则打印日志
		log.Println(err)
	}
	return errors.New("尝试了所有的服务商，但是都发送失败了")
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}
