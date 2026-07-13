package router

import (
	"github.com/linbaozhong/gentity/pkg/ack/iris"
	_ "tui/internal/handler"
)

func Init(port string) *ack.Server {
	app := ack.NewApplication("tui", "0.1")

	v1 := ack.NewParty(app, "/v1")
	// 注册路由
	ack.RegisterRouter(v1)

	if port[0] != ':' {
		port = ":" + port
	}
	return ack.NewServer(app, port)
}
