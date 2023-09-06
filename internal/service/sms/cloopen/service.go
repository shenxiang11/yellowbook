package cloopen

import (
	"context"
	"errors"
	"github.com/shenxiang11/go-sms-sdk/cloopen"
	"log"
	"strings"
	"yellowbook/config"
	"yellowbook/internal/service/sms"
)

var ErrSMSSendFailed = errors.New("发送异常")

type Service struct {
	appId  string
	client cloopen.IClient
}

func NewService(client cloopen.IClient) sms.Service {
	return &Service{
		appId:  config.Conf.Cloopen.AppId,
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
		return ErrSMSSendFailed
	}

	return nil
}
