package failover

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project_go/webbook/internal/service/sms"
	smsmocks "project_go/webbook/internal/service/sms/mocks"
	"testing"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		wantErr error
	}{
		{
			name: "最后一个发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := []sms.Service{svc1, svc2}
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svcs
			},
			wantErr: nil,
		},
		{
			name: "都发送失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svcs := []sms.Service{svc1, svc2}
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				return svcs
			},
			wantErr: errors.New("尝试了所有的服务商，但是都发送失败了"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svcs := tc.mock(ctrl)
			failOverSmsService := NewFailOverSMSService(svcs)
			err := failOverSmsService.Send(context.Background(), "12", []string{"122"})
			assert.Equal(t, tc.wantErr, err)
		})
	}

}
