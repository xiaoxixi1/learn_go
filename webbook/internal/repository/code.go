package repository

import (
	"context"
	"project_go/webbook/internal/repository/cache"
)

var (
	VerifyTooManyError = cache.VerifyTooManyError
	SendTooManyError   = cache.SendTooManyError
)

type CodeRepo struct {
	codeCache *cache.CodeCache
}

func NewCodeRepo(codeCache *cache.CodeCache) *CodeRepo {
	return &CodeRepo{
		codeCache: codeCache,
	}
}

func (c *CodeRepo) Set(cxt context.Context, biz, phone, code string) error {
	return c.codeCache.Set(cxt, biz, phone, code)
}
func (c *CodeRepo) Verify(cxt context.Context, biz, phone, code string) (bool, error) {
	return c.codeCache.Verify(cxt, biz, phone, code)
}
