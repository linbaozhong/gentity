package router

import (
	_ "abc/internal/handler"
	"github.com/linbaozhong/gentity/pkg/api"
)

func Init() api.Application {
	_app := api.NewApplication("abc", "0.1")

	_app.Use(api.Recovery())
	_app.Use(api.Logger())

	_v1 := api.NewParty(_app, "/v1")
	// 注册路由
	api.RegisterRouter(_v1)
	return _app
}
