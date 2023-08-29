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
	engine := InitWebServer()

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
