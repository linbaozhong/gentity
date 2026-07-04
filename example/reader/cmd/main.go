package main

import (
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"os"
	"os/signal"
	"reader/internal/router"
	"syscall"
	"time"
)

var (
	_ = app.Context
)

func main() {
	port := ":8080"
	// 命令行指定端口
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	if port[0] != ':' {
		port = ":" + port
	}

	log.Register(false)
	log.Info(fmt.Sprintf("%s %s %s 服务已开启", "reader Api", "0.1", port))
	// 启动API服务
	_app := router.Init()

	app.Launch()

	if err := _app.Listen(port); err != nil {
		log.Error(err)
	}

	idleConnsClosed := make(chan os.Signal, 1)
	signal.Notify(idleConnsClosed, syscall.SIGINT, syscall.SIGTERM)
	// 优雅地关闭
	<-idleConnsClosed
	close(idleConnsClosed)

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// close all hosts.
	if err := _app.Shutdown(ctx); err != nil {
		log.Info("Server Shutdown:", err)
	}
	closing("reader Api", "0.1", port)
}
func closing(name, ver, addr string) {
	app.Close()
	log.Fatal(fmt.Sprintf("%s %s %s 服务已关闭", name, ver, addr))
}
