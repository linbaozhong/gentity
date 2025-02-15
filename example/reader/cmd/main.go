package main

import (
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/gentity/pkg/api/broker"
	"github.com/linbaozhong/gentity/pkg/log"
	"os"
	"reader/internal/router"
	"reader/internal/service"
	"time"
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

	log.RegisterLogger(broker.Context, false)
	log.Info(fmt.Sprintf("%s %s %s 服务已开启", "reader Api", "0.1", port))
	// 启动API服务
	app := router.Init()

	idleConnsClosed := make(chan struct{})
	api.OnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// close all hosts.
		app.Shutdown(ctx)
		closing("reader Api", "0.1", port)
		close(idleConnsClosed)
	})

	service.Open(api.AppContext)

	if err := app.Listen(port); err != nil {
		log.Error(err)
	}

	// 优雅地关闭
	<-idleConnsClosed
}
func closing(name, ver, addr string) {
	broker.Close()
	log.Fatal(fmt.Sprintf("%s %s %s 服务已关闭", name, ver, addr))
}
