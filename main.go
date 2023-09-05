package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"yellowbook/config"
)

func main() {
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
