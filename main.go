package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/web/middleware"
)
import "gorm.io/driver/mysql"

func main() {
	db := initDB()

	server := initWebServer()

	u := InitUserHandler(db)
	u.RegisterRoutes(server.Group("/users"))

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		MaxAge:           2 * time.Minute,
	}))

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("yellow-id", store))

	server.Use(middleware.NewLoinMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())

	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13306)/yellowbook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
