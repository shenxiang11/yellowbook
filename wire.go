//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"yellowbook/internal/repository"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/service"
	"yellowbook/internal/service/sms/cloopen"
	"yellowbook/internal/web"
)

func InitUserHandler(smsAppId string) *web.UserHandler {
	wire.Build(
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
		cache.NewUserCache,
		service.NewCodeService,
		repository.NewCodeRepository,
		cache.NewCodeCache,
		cloopen.NewService,

		initDB,
		initRedis,
		initCloopen,
	)
	return &web.UserHandler{}
}
