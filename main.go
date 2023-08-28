package main

import (
	"context"
	"github.com/cloopen/go-sms-sdk/cloopen"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yellowbook/config"
	"yellowbook/internal/repository/dao"
	"yellowbook/internal/web/middleware"
)
import "gorm.io/driver/mysql"

func main() {
	engine := initWebServer()

	u := InitUserHandler()
	u.RegisterRoutes(engine.Group("/users"))

	server := &http.Server{
		Addr:    config.Conf.Web.Port,
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

	server.Use(
		middleware.NewLoinMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			Build(),
	)

	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Conf.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Conf.Redis.Addr,
	})
	return redisClient
}

func initCloopen() *cloopen.Client {
	cfg := cloopen.DefaultConfig().
		WithAPIAccount("8aaf07087fe90a32017ff389d6ac01bb").
		WithAPIToken("a1c23065a7d847c384d719ad240f6384")

	client := cloopen.NewJsonClient(cfg)

	return client
}
