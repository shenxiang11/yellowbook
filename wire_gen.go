// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"yellowbook/internal/manage"
	"yellowbook/internal/repository"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/cache/ristretto"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/service"
	"yellowbook/internal/web"
	"yellowbook/ioc"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	db := ioc.InitDB()
	userDao := dao.NewUserDAO(db)
	cmdable := ioc.InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDao, userCache)
	iUserService := service.NewUserService(userRepository)
	ristrettoCache := ioc.InitRistretto()
	codeCache := ristretto.NewCodeCache(ristrettoCache)
	codeRepository := repository.NewCodeRepository(codeCache)
	client := ioc.InitCloopen()
	smsService := ioc.InitSMSService(client)
	codeService := service.NewCodeService(codeRepository, smsService)
	iService := ioc.InitGithub()
	ijwtGenerator := ioc.InitJWT()
	userHandler := web.NewUserHandler(iUserService, codeService, iService, ijwtGenerator)
	ossIService := ioc.InitOss()
	iResourceDao := dao.NewResourceDAO(db)
	iResourceRepository := repository.NewResourceRepository(iResourceDao)
	logger := ioc.InitLogger()
	iResourceService := service.NewResourceService(ossIService, iResourceRepository, logger)
	resourceHandler := web.NewResourceHandler(iResourceService)
	iArticleDAO := dao.NewArticleDAO(db)
	iArticleRepository := repository.NewArticleRepository(iArticleDAO)
	iArticleService := service.NewArticleService(iArticleRepository, logger)
	articleHandler := web.NewArticleHandler(iArticleService)
	engine := ioc.InitWebServer(userHandler, resourceHandler, articleHandler, logger)
	return engine
}

func InitManageServer() *gin.Engine {
	db := ioc.InitDB()
	userDao := dao.NewUserDAO(db)
	cmdable := ioc.InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDao, userCache)
	iUserService := service.NewUserService(userRepository)
	userHandler := manage.NewUserHandler(iUserService)
	engine := ioc.InitManageServer(userHandler)
	return engine
}

func InitSpider() *ioc.Spider {
	db := ioc.InitDB()
	iArticleDAO := dao.NewArticleDAO(db)
	iArticleRepository := repository.NewArticleRepository(iArticleDAO)
	logger := ioc.InitLogger()
	iArticleService := service.NewArticleService(iArticleRepository, logger)
	spider := ioc.NewSpider(iArticleService)
	return spider
}
