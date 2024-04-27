package repository

import (
	"context"
	"project_go/webbook/internal/repository/cache"
)

var (
	VerifyTooManyError = cache.VerifyTooManyError
	SendTooManyError   = cache.SendTooManyError
)

type CodeRepo interface {
	Set(cxt context.Context, biz, phone, code string) error
	Verify(cxt context.Context, biz, phone, code string) (bool, error)
}

type CachedCodeRepo struct {
	codeCache cache.CodeCache
}

func NewCodeRepo(codeCache cache.CodeCache) CodeRepo {
	return &CachedCodeRepo{
		codeCache: codeCache,
	}
}

func (c *CachedCodeRepo) Set(cxt context.Context, biz, phone, code string) error {
	return c.codeCache.Set(cxt, biz, phone, code)
}
func (c *CachedCodeRepo) Verify(cxt context.Context, biz, phone, code string) (bool, error) {
	return c.codeCache.Verify(cxt, biz, phone, code)
}
