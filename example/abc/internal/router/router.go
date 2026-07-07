package router

import (
	_ "github.com/linbaozhong/gentity/example/abc/internal/handler"
	api "github.com/linbaozhong/gentity/pkg/api/iris"
)

func Init(port string) *api.Server {
	_app := api.NewApplication("abc", "0.1")

	_app.Use(api.Recovery())
	_app.Use(api.Logger())

	_v1 := api.NewParty(_app, "/v1")
	// 注册路由
	api.RegisterRouter(_v1)

	if port[0] != ':' {
		port = ":" + port
	}
	return api.NewServer(_app, port)
}
