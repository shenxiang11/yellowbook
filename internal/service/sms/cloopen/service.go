package cloopen

import (
	"context"
	"errors"
	"github.com/cloopen/go-sms-sdk/cloopen"
	"log"
	"strings"
	"yellowbook/internal/service/sms"
)

type Service struct {
	appId  string
	client *cloopen.Client
}

func NewService(client *cloopen.Client, appId string) *Service {
	return &Service{
		appId:  appId,
		client: client,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args sms.NamedArgSlice, to ...string) error {
	input := &cloopen.SendRequest{
		AppId:      s.appId,
		To:         strings.Join(to, ","),
		TemplateId: tpl,
		Datas:      sms.ConvertToStrSlice(args),
	}

	resp, err := s.client.SMS().Send(input)
	if err != nil {
		return err
	}

	if resp.StatusCode != "000000" {
		log.Println(resp)
		return errors.New(resp.StatusMsg)
	}

	return nil
}
