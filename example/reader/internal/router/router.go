package router

import (
	"github.com/linbaozhong/gentity/pkg/api"
	_ "reader/internal/handler"
)

func Init() api.Application {
	app := api.NewApplication("reader", "0.1")

	app.Use(api.Recovery())
	app.Use(api.Logger())

	v1 := api.NewParty(app, "/v1")
	// 注册路由
	l := len(api.Instances)
	for i := 0; i < l; i++ {
		if m, ok := api.Instances[i].(api.IRegisterRouter); ok {
			m.RegisterRouter(v1)
		}
	}
	return app
}
