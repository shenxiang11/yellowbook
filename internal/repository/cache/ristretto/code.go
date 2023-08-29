package ristretto

import (
	"context"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"log"
	"sync"
	"time"
	"yellowbook/internal/repository/cache"
)

type CodeCache struct {
	cache *ristretto.Cache
	ttl   time.Duration
	mux   sync.Mutex
}

func NewCodeCache(cache *ristretto.Cache) cache.CodeCache {
	return &CodeCache{
		cache: cache,
		ttl:   time.Minute * 10,
	}
}

type CodeCacheData struct {
	code  string
	count int
}

func (c *CodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	// 先查询 ttl，如果没有过期时间，则是系统错误
	// 如果 ttl 大于 ttl-1 分钟，说明一分中内发送多次
	// 其他 则允许发送

	c.mux.Lock()
	defer c.mux.Unlock()

	key := c.key(biz, phone)

	ttl, ok := c.cache.GetTTL(key)
	if ok && ttl > c.ttl-time.Minute {
		return cache.ErrCodeSendTooMany
	}

	ok = c.cache.SetWithTTL(
		key,
		CodeCacheData{
			code:  code,
			count: 3,
		},
		0,
		c.ttl,
	)
	if !ok {
		// 不知道怎么，设置失败了，打日志排查呗
		log.Println("设置失败")
		return cache.ErrUnknown
	}

	return nil
}

func (c *CodeCache) Verify(ctx context.Context, biz string, phone string, code string) error {
	// 取出来，查看次数，如果小于等于 0 次，则是在攻击
	// 如果 code 匹配，则正确
	// 如果 code 不匹配，次数减少 1 次

	c.mux.Lock()
	defer c.mux.Unlock()

	key := c.key(biz, phone)

	data, ok := c.cache.Get(key)
	if !ok {
		return cache.ErrUnknown
	}

	d, ok := data.(CodeCacheData)
	if !ok {
		return cache.ErrUnknown
	}

	if d.count <= 0 {
		return cache.ErrCodeVerifyTooManyTimes
	} else {
		if d.code == code {
			return nil
		} else {
			return cache.ErrCodeVerifyFailed
		}
	}
}

func (c *CodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
