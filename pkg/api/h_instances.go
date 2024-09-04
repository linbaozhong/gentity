package api

var (
	Instances = make([]interface{}, 0)
)

type IRegisterRouter interface {
	RegisterRouter(party Party)
}

type base struct{}
