package ioc

import (
	"project_go/webbook/internal/service/sms"
	"project_go/webbook/internal/service/sms/localmodel"
)

func InitSmsSendService() sms.Service {
	// 没有腾讯云的APP ID，使用本地的
	return localmodel.NewService()
}
