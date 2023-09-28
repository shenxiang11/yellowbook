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
		web.NewResourceHandler,
		web.NewUserHandler,
		web.NewArticleHandler,

		service.NewUserService,
		service.NewResourceService,
		service.NewArticleService,
		service.NewCodeService,

		repository.NewCachedUserRepository,
		repository.NewResourceRepository,
		repository.NewArticleRepository,
		repository.NewCodeRepository,

		dao.NewResourceDAO,
		dao.NewUserDAO,
		dao.NewArticleDAO,

		cache.NewUserCache,
		ristretto.NewCodeCache,

		ioc.InitOss,
		ioc.InitRistretto,
		ioc.InitWebServer,
		ioc.InitSMSService,
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitCloopen,
		ioc.InitJWT,
		ioc.InitGithub,
		ioc.InitLogger,
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

func InitSpider() *ioc.Spider {
	wire.Build(
		ioc.InitLogger,
		ioc.InitDB,
		dao.NewArticleDAO,
		repository.NewArticleRepository,
		service.NewArticleService,
		ioc.NewSpider,
	)
	return &ioc.Spider{}
}
