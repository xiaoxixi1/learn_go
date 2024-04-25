package repository

import (
	"context"
	"project_go/webbook/internal/repository/cache"
)

var (
	VerifyTooManyError = cache.VerifyTooManyError
)

type CodeRepo struct {
	codeCache *cache.Codecache
}

func NewCodeRepo(codeCache cache.Codecache) *CodeRepo {
	return &CodeRepo{
		codeCache: &codeCache,
	}
}

func (c *CodeRepo) Set(cxt context.Context, biz, phone, code string) error {
	return c.codeCache.Set(cxt, biz, phone, code)
}
func (c *CodeRepo) Verify(cxt context.Context, biz, phone, code string) (bool, error) {
	return c.codeCache.Verify(cxt, biz, phone, code)
}
