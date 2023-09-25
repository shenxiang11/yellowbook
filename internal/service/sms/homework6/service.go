package homework6

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/shenxiang11/zippo/slice"
	"time"
	"yellowbook/internal/service/sms"
)

// 作业，可以多服务商，失败转异步重试的短信服务

var ErrNoAvailable = errors.New("无可用服务商")
var ErrUnknown = errors.New("未知错误")

type Service struct {
	svcs    map[string]sms.Service // 需要多个短信服务
	client  redis.Cmdable          // 需要一个 redis，帮忙存储被踢掉的服务商
	keys    []string
	banTime time.Duration
}

func NewService(svcs map[string]sms.Service, client redis.Cmdable, banTime time.Duration) sms.Service {
	keys := []string{}

	for key, _ := range svcs {
		keys = append(keys, key)
	}

	return &Service{
		svcs:    svcs,
		client:  client,
		banTime: banTime,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args sms.NamedArgSlice, to ...string) error {
	banned, err := s.client.Keys(ctx, "banned_sms").Result()
	if err != nil {
		return err
	}

	availableSms := slice.Filter(s.keys, func(kel string, index int) bool {
		return !slice.Some(banned, func(bel string, index int) bool {
			return kel == bel
		})
	})

	if len(availableSms) == 0 {
		return ErrNoAvailable
	}

	// 选一个服务发

	choose, ok := s.svcs[availableSms[0]]
	if !ok {
		return ErrUnknown
	}

	err = choose.Send(ctx, tpl, args, to...)
	if err != nil {
		// 把它添加到黑名单中
		//s.client.SetEx()

		// 转异步重试
		return err
	}

	return nil
}
