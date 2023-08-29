package ioc

import (
	"github.com/cloopen/go-sms-sdk/cloopen"
	"yellowbook/internal/service/sms"
	cloopen2 "yellowbook/internal/service/sms/cloopen"
)

func InitSMSService(c *cloopen.Client) sms.Service {
	return cloopen2.NewService(c)
}

func InitCloopen() *cloopen.Client {
	cfg := cloopen.DefaultConfig().
		WithAPIAccount("8aaf07087fe90a32017ff389d6ac01bb").
		WithAPIToken("a1c23065a7d847c384d719ad240f6384")

	client := cloopen.NewJsonClient(cfg)

	return client
}
