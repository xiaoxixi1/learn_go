package localmodel

import (
	"context"
	"fmt"
)

type Service struct {
}

func (s Service) Send(cxt context.Context, tplId string, args []string, number ...string) error {
	fmt.Printf("发送一条验证码短信：%s", args[0])
	return nil
}

func NewService() *Service {
	return &Service{}
}
