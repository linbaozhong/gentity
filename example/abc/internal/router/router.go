package router

import (
	_ "github.com/linbaozhong/gentity/example/abc/internal/handler"
	"github.com/linbaozhong/gentity/pkg/ack/iris"
)

func Init(port string) *ack.Server {
	_app := ack.NewApplication("abc", "0.1")

	_v1 := ack.NewParty(_app, "/v1")

	// 注册路由
	ack.RegisterRouter(_v1)

	if port[0] != ':' {
		port = ":" + port
	}
	return ack.NewServer(_app, port)
}
