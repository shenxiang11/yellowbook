package cache

import (
	"context"
	"errors"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证错误次数太多")
	ErrCodeVerifyFailed       = errors.New("验证码错误")
	ErrUnknown                = errors.New("unknown")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) error
}
