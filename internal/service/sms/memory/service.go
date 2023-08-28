package memory

import (
	"context"
	"fmt"
	"yellowbook/internal/service/sms"
)

type Service struct {
}

func NewService() sms.Service {
	return &Service{}
}

func (s Service) Send(ctx context.Context, tpl string, args sms.NamedArgSlice, to ...string) error {
	fmt.Println(args, to)
	return nil
}
