package service

import (
	"context"
	"fmt"
	"math/rand"
	"yellowbook/internal/repository"
	"yellowbook/internal/service/sms"
)

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeVerifyFailed       = repository.ErrCodeVerifyFailed
	ErrUnknown                = repository.ErrUnknown
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, code string) error
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	return svc.smsSvc.Send(ctx, "1", []sms.NamedArg{
		{
			Name: "1",
			Val:  code,
		},
		{
			Name: "2",
			Val:  "10",
		},
	}, phone)
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, code string) error {
	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *codeService) generateCode() string {
	num := rand.Intn(10000)
	return fmt.Sprintf("%04d", num)
}
