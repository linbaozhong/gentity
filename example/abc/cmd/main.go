package main

import (
	"abc/internal/router"
	"abc/internal/service"
	"context"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/serverpush"
	"os"
	"time"
)

var (
	_ = app.Context
)

func main() {
	_port := ":8080"
	// 命令行指定端口
	if len(os.Args) > 1 {
		_port = os.Args[1]
	}
	if _port[0] != ':' {
		_port = ":" + _port
	}

	log.Register(false)
	log.Info(fmt.Sprintf("%s %s %s 服务开启中...", "abc Api", "0.1", _port))

	// 启动前置服务
	Prepare()
	// 启动路由服务
	_service := router.Init()

	_idleConnsClosed := make(chan struct{})
	api.OnInterrupt(func() {
		_timeout := 5 * time.Second
		_ctx, _cancel := context.WithTimeout(context.Background(), _timeout)
		defer _cancel()
		// close all hosts.
		_service.Shutdown(_ctx)
		// 关闭其他服务
		Finished()
		log.Fatal(fmt.Sprintf("%s %s %s 服务已关闭", "abc Api", "0.1", _port))
		close(_idleConnsClosed)
	})

	if e := _service.Listen(_port); e != nil {
		log.Error(e)
	}

	// 优雅地关闭
	<-_idleConnsClosed
}

// Prepare 系统启动所需要的必须服务
func Prepare() {
	// 连接数据库
	service.Connected()
	// 启动推送服务
	serverpush.Start(
		serverpush.WithAutoStream(true),
		serverpush.WithAutoReplay(true),
	)
}

// Finished 系统关闭后结束所有其他的服务
func Finished() {
	// 关闭推送服务
	serverpush.Close()
	app.Close()
}
