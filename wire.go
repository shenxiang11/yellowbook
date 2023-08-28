//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"yellowbook/internal/repository"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/service"
	"yellowbook/internal/service/sms/memory"
	"yellowbook/internal/web"
)

func InitUserHandler() *web.UserHandler {
	wire.Build(
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
		cache.NewUserCache,
		service.NewCodeService,
		repository.NewCodeRepository,
		cache.NewCodeCache,
		memory.NewService,

		initDB,
		initRedis,
	)
	return &web.UserHandler{}
}
