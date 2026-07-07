package main

import (
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/example/reader/internal/router"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
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

	log.Register(false)
	log.Info(fmt.Sprintf("%s %s %s 服务已开启", "reader Api", "0.1", port))
	// 启动API服务
	srv := router.Init(port)

	app.Launch()

	// 监听系统信号，实现优雅停止
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Fatal(fmt.Sprintf("%s %s %s 服务已关闭", "reader Api", "0.1", port))

		// 给正在处理的请求 5 秒钟的时间完成
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 关闭其他服务
		closing("reader Api", "0.1", port)
		// 优雅关闭
		srv.Shutdown(ctx)
	}()

	// 启动服务
	if err := srv.Run(); err != nil {
		log.Error("server error", err)
	}

}
func closing(name, ver, addr string) {
	app.Close()
	log.Fatal(fmt.Sprintf("%s %s %s 服务已关闭", name, ver, addr))
}
