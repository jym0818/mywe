package repository

import (
	"context"

	"github.com/jym0818/mywe/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = cache.ErrUnknownForCode
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeRepository struct {
	cache cache.CodeCache
}

func (c *codeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *codeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)
}

func NewcodeRepository(cache cache.CodeCache) CodeRepository {
	return &codeRepository{
		cache: cache,
	}
}
