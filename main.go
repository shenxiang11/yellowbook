package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/web/middleware"
)
import "gorm.io/driver/mysql"

func main() {
	db := initDB()

	engine := initWebServer()

	u := InitUserHandler(db)
	u.RegisterRoutes(engine.Group("/users"))

	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			// 优雅退出不能写 panic
			// panic(err)
			log.Println("Server err: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Jwt-Token"},
		MaxAge:           2 * time.Minute,
	}))

	store, err := redis.NewStore(16, "tcp", "localhost:16379", "", []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	if err != nil {
		panic(err)
	}

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
