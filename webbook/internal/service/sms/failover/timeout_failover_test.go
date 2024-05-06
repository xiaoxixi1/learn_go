package failover

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project_go/webbook/internal/service/sms"
	smsmocks "project_go/webbook/internal/service/sms/mocks"
	"testing"
)

func TestTimeOutFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		threshold int32
		idx       int32
		cnt       int32
		wantErr   error
		wantIdx   int32
		wantCnt   int32
	}{
		{
			name: "没有达到超时阈值,发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := []sms.Service{svc1, svc2}
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svcs
			},
			threshold: 15,
			cnt:       12,
			idx:       0,
			wantCnt:   0,
			wantIdx:   0,
		},
		{
			name: "达到超时阈值切换svc2,发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := []sms.Service{svc1, svc2}
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svcs
			},
			idx:       0,
			cnt:       12,
			threshold: 12,
			wantIdx:   1,
			wantCnt:   0,
		},
		{
			name: "没有达到超时阈值，发送失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := []sms.Service{svc1, svc2}
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return svcs
			},
			threshold: 15,
			cnt:       12,
			idx:       0,
			wantErr:   context.DeadlineExceeded,
			wantCnt:   13,
			wantIdx:   0,
		},
		{
			name: "达到超时阈值，发送失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := []sms.Service{svc1, svc2}
				//svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return svcs
			},
			threshold: 12,
			cnt:       12,
			idx:       0,
			wantErr:   context.DeadlineExceeded,
			wantCnt:   1,
			wantIdx:   1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svcs := tc.mock(ctrl)
			tfs := NewTimeOutFailOverSMSService(svcs, tc.threshold)
			tfs.idx = tc.idx
			tfs.cnt = tc.cnt
			err := tfs.Send(context.Background(), "1", []string{"1"})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, tfs.idx)
			assert.Equal(t, tc.wantCnt, tfs.cnt)
		})
	}

}
