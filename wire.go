//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"yellowbook/internal/repository"
	"yellowbook/internal/repository/cache"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/service"
	"yellowbook/internal/web"
)

func InitUserHandler(db *gorm.DB, redisCmd redis.Cmdable) *web.UserHandler {
	wire.Build(
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
		cache.NewUserCache,
	)
	return &web.UserHandler{}
}
