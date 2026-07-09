package router

import (
	_ "github.com/linbaozhong/gentity/example/reader/internal/handler"
	"github.com/linbaozhong/gentity/pkg/ack/gin"
)

func Init(port string) *ack.Server {
	app := ack.NewApplication("reader", "0.1")

	v1 := ack.NewParty(app, "/v1")
	// 注册路由
	ack.RegisterRouter(v1)

	if port[0] != ':' {
		port = ":" + port
	}
	return ack.NewServer(app, port)
}
