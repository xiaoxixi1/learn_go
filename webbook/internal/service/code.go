package service

import (
	"context"
	"fmt"
	"math/rand"
	"project_go/webbook/internal/repository"
	"project_go/webbook/internal/service/sms"
)

var tplId = "1877556"
var ErrCodeSendTooMany = repository.SendTooManyError

type CodeService interface {
	SendCode(cxt context.Context, biz, phone string) error
	Verify(cxt context.Context, biz string, phone string, inputCode string) (bool, error)
}
type codeService struct {
	repo repository.CodeRepo
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepo, sms sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  sms,
	}
}

func (cs *codeService) SendCode(cxt context.Context, biz, phone string) error {
	code := cs.generate()
	fmt.Printf("验证码：%s", code)
	err := cs.repo.Set(cxt, biz, phone, code)
	if err != nil {
		return err
	}
	return cs.sms.Send(cxt, tplId, []string{code}, phone)

}

func (cs *codeService) Verify(cxt context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := cs.repo.Verify(cxt, biz, phone, inputCode)
	if err == repository.VerifyTooManyError {
		// 相当于对外面屏蔽了验证次数过多的错误，只告诉调用着这个不错
		// 但是对于业务内部，此处可以埋点，告警
		return false, nil
	}
	return ok, err
}

func (cs *codeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)

}
