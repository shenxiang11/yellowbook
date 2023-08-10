//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"yellowbook/internal/repository"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/service"
	"yellowbook/internal/web"
)

func InitUserHandler(db *gorm.DB) *web.UserHandler {
	wire.Build(
		web.NewUserHandler,
		service.NewUserService,
		repository.NewUserRepository,
		dao.NewUserDAO,
	)
	return &web.UserHandler{}
}
