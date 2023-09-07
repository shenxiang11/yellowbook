package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"yellowbook/internal/manage"
)

func InitManageServer(userHandler *manage.UserHandler) *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Jwt-Token"},
		MaxAge:           2 * time.Minute,
	}))

	//server.Use(
	//	middleware.NewLoginMiddlewareBuilder().
	//		Build(),
	//)

	userHandler.RegisterRoutes(server.Group("/users"))

	return server
}
