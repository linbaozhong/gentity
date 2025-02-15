package api

import (
	"github.com/linbaozhong/gentity/pkg/app"
	"strconv"
	"time"
)

var (
	Instances = make([]any, 0)
)

type IRegisterRouter interface {
	RegisterRouter(party Party)
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
