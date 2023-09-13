//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"yellowbook/internal/manage"
	"yellowbook/internal/repository"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/cache/ristretto"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/service"
	"yellowbook/internal/web"
	"yellowbook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		web.NewUserHandler,
		service.NewUserService,
		repository.NewCachedUserRepository,
		dao.NewUserDAO,
		cache.NewUserCache,
		service.NewCodeService,
		repository.NewCodeRepository,
		ristretto.NewCodeCache,

		ioc.InitRistretto,
		ioc.InitWebServer,
		ioc.InitSMSService,
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitCloopen,
		ioc.InitJWT,
		ioc.InitGithub,
	)
	return new(gin.Engine)
}

func InitManageServer() *gin.Engine {
	wire.Build(
		manage.NewUserHandler,
		service.NewUserService,
		repository.NewCachedUserRepository,
		dao.NewUserDAO,
		cache.NewUserCache,

		ioc.InitManageServer,
		ioc.InitDB,
		ioc.InitRedis,
	)
	return new(gin.Engine)
}
