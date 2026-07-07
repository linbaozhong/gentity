package router

import (
	_ "github.com/linbaozhong/gentity/example/reader/internal/handler"
	api "github.com/linbaozhong/gentity/pkg/api/gin"
)

func Init(port string) *api.Server {
	app := api.NewApplication("reader", "0.1")

	app.Use(api.Recovery())
	app.Use(api.Logger())

	v1 := api.NewParty(app, "/v1")
	// 注册路由
	api.RegisterRouter(v1)

	if port[0] != ':' {
		port = ":" + port
	}
	return api.NewServer(app, port)
}
