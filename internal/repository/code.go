package repository

import (
	"context"

	"github.com/ra1n6ow/webook/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepositoryer interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CodeRepository struct {
	cache cache.CodeCacher
}

func NewCodeRepository(cache cache.CodeCacher) CodeRepositoryer {
	return &CodeRepository{
		cache: cache,
	}
}

func (r *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return r.cache.Set(ctx, biz, phone, code)
}

// bool 表示是否验证成功, 业务向
// error 表示验证失败的原因, 代码向
func (r *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return r.cache.Verify(ctx, biz, phone, code)
}
