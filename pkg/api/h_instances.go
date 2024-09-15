package api

var (
	Instances = make([]any, 0)
)

type IRegisterRouter interface {
	RegisterRouter(party Party)
}

type base struct{}
