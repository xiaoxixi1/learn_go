package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
)

var (
	SendTooManyError   = errors.New("发送太频繁")
	VerifyTooManyError = errors.New("验证码已失效")
)

type CodeCache interface {
	Set(cxt context.Context, biz, phone, code string) error
	Verify(cxt context.Context, biz, phone, code string) (bool, error)
}

func Key(biz string, phone string) string {
	fmt.Printf("phone_code:%s:%s", biz, phone)
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
