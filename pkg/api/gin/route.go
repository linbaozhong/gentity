package api

var (
	routes = make([]any, 0)
)

type IRegisterRoute interface {
	RegisterRoute(party Party)
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
