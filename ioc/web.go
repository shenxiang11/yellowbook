package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"yellowbook/internal/web"
	"yellowbook/internal/web/middleware"
	"yellowbook/pkg/logger"
)

func InitWebServer(
	userHandler *web.UserHandler,
	resourceHandler *web.ResourceHandler,
	articleHandler *web.ArticleHandler,
	l logger.Logger,
) *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Jwt-Token"},
		MaxAge:           2 * time.Minute,
	}))

	server.Use(
		middleware.NewLoggerMiddlewareBuilder(func(ctx context.Context, al *middleware.AccessLog) {
			l.Debug("HTTP 请求", logger.Field{Key: "Access Log: ", Value: al})
		}).
			Build(),
	)

	server.Use(
		middleware.NewLoginMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/users/github/oauth").
			IgnorePaths("/users/github/authorize").
			IgnorePaths("/users/version").
			Build(),
	)

	userHandler.RegisterRoutes(server.Group("/users"))
	resourceHandler.RegisterRoutes(server.Group("/resources"))
	articleHandler.RegisterRoutes(server.Group("/articles"))

	return server
}
