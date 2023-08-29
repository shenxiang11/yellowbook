package redis

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"yellowbook/internal/repository/cache"
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) cache.CodeCache {
	return &CodeCache{client: client}
}

func (c *CodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return cache.ErrCodeSendTooMany
	default:
		log.Println("未知错误：", err)
		return cache.ErrUnknown
	}
}

func (c *CodeCache) Verify(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return cache.ErrCodeVerifyTooManyTimes
	case -2:
		return cache.ErrCodeVerifyFailed
	default:
		log.Println("未知错误：", err)
		return cache.ErrUnknown
	}
}

func (c *CodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
