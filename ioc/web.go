package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"yellowbook/internal/web"
	"yellowbook/internal/web/middleware"
)

func InitWebServer(userHandler *web.UserHandler) *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Jwt-Token"},
		MaxAge:           2 * time.Minute,
	}))

	server.Use(
		middleware.NewLoinMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			Build(),
	)

	userHandler.RegisterRoutes(server.Group("/users"))

	return server
}
