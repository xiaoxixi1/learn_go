package ratelimit

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project_go/webbook/internal/service/sms"
	smsmocks "project_go/webbook/internal/service/sms/mocks"
	"project_go/webbook/pkg/limiter"
	limitermocks "project_go/webbook/pkg/limiter/mocks"
	"testing"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		// 测试用例跟输入没有关系，所以可以写死
		wantErr error
	}{
		{
			name: "没有限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limit := limitermocks.NewMockLimiter(ctrl)
				limit.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, limit
			},
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limit := limitermocks.NewMockLimiter(ctrl)
				limit.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return svc, limit
			},
			wantErr: errLimited,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limit := limitermocks.NewMockLimiter(ctrl)
				limit.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("redis限流错误"))
				return svc, limit
			},
			wantErr: errors.New("redis限流错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, limit := tc.mock(ctrl)
			rateLimitSmsService := NewRateLimitSMSService(smsSvc, limit, "test")
			err := rateLimitSmsService.Send(context.Background(), "11", []string{"ss"})
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
