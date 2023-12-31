package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"yellowbook/internal/manage"
)

func InitManageServer(userHandler *manage.UserHandler, articleHandler *manage.ArticleHandler) *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Yellow-Book-Timezone"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Jwt-Token"},
		MaxAge:           2 * time.Minute,
	}))

	//server.Use(
	//	middleware.NewLoginMiddlewareBuilder().
	//		Build(),
	//)

	userHandler.RegisterRoutes(server.Group("/users"))
	articleHandler.RegisterRoutes(server.Group("/articles"))

	return server
}
