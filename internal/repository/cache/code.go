package cache

import (
	"context"
	"errors"
	"fmt"

	_ "embed"

	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

type CodeCacher interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) CodeCacher {
	return &CodeCache{
		cmd: cmd,
	}
}

func (c *CodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone:code:%s:%s", biz, phone)
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// 因为先查询再操作会有并发问题
	// 通过 lua 脚本设置验证码, 是原子操作
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *CodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 因为先查询再操作会有并发问题
	// 通过 lua 脚本设置验证码, 是原子操作
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}
