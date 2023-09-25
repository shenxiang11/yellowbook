package main

import (
	"context"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yellowbook/config"
)

func main() {
	initViper()

	webEngine := InitWebServer()
	manageEngine := InitManageServer()

	webServer := &http.Server{
		Addr:    config.Conf.Web.Port,
		Handler: webEngine,
	}

	manageServer := &http.Server{
		Addr:    config.Conf.Manage.Port,
		Handler: manageEngine,
	}

	go func() {
		err := webServer.ListenAndServe()
		if err != nil {
			// 优雅退出不能写 panic
			// panic(err)
			log.Println("Server err: ", err)
		}
	}()

	go func() {
		err := manageServer.ListenAndServe()
		if err != nil {
			// 优雅退出不能写 panic
			// panic(err)
			log.Println("Server err: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	if err := webServer.Shutdown(context.Background()); err != nil {
		log.Fatal("web server shutdown failed:", err)
	}

	if err := manageServer.Shutdown(context.Background()); err != nil {
		log.Fatal("manage server shutdown failed:", err)
	}
}

func initViper() {
	err := viper.AddRemoteProvider("consul", config.Conf.Consul.DSN, config.Conf.Consul.Key)
	if err != nil {
		return
	}

	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}

	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			// 注意：修改远程配置会有 5s 的延迟才能生效
			time.Sleep(time.Second * 5)
			err = viper.ReadRemoteConfig()
			if err != nil {
				log.Println("更新配置失败！！")
			}
		}
	}()
}
