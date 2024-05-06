package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"project_go/webbook/internal/service/sms"
)

/**
  短信服务商，提供安全性
  一般短信服务，会先申请一个模板，提交自己业务的一些信息，比如流量，等等
  然后内部审核通过之后，一般会生成一个静态的token
  我们将返回给业务方的token中，封装一些跟业务方有关的信息
  在使用时，直接将各种数据从token中解密出来
*/

/*
*

	短信发送内部鉴权
*/
type SMSService struct {
	svc sms.Service
	key []byte
}

func (s *SMSService) Send(cxt context.Context, tplToken string, args []string, number ...string) error {
	var claim SMSClaim
	_, err := jwt.ParseWithClaims(tplToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(cxt, claim.tpl, args, number...)
}

type SMSClaim struct {
	jwt.RegisteredClaims
	tpl string
	// 可以额外加字段
}
