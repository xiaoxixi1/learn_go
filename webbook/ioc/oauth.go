package ioc

import (
	"project_go/webbook/internal/service/oauth2/wechat"
)

func InitWeChatOAuthService() wechat.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("找不到appID的环境变量")
	//}
	appId := "123"
	secretId := "456"
	return wechat.NewService(appId, secretId)
}
