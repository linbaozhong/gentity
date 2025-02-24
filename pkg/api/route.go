package api

import (
	"github.com/linbaozhong/gentity/pkg/app"
	"strconv"
	"time"
)

var (
	routes = make([]any, 0)
)

type IRegisterRoute interface {
	RegisterRoute(party Party)
}

func Initiate(ctx Context, arg any) {
	if id := ctx.GetHeader(app.OperationID); len(id) == 0 {
		ctx.SetID(strconv.FormatInt(time.Now().UnixMilli(), 10))
	} else {
		ctx.SetID(id)
	}

	if ier, ok := arg.(Initializer); ok {
		ier.Init()
	}
}
func InitiateX(ctx Context, arg any) {
	if ier, ok := arg.(Initializer); ok {
		ier.Init()
	}
}

// RegisterRoute 注册路由
func RegisterRoute(r IRegisterRoute) {
	routes = append(routes, r)
}

// RegisterRouter 注册路由器
func RegisterRouter(group Party) {
	_l := len(routes)
	for i := 0; i < _l; i++ {
		if m, ok := routes[i].(IRegisterRoute); ok {
			m.RegisterRoute(group)
		}
	}
}
