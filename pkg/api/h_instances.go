package api

import (
	"github.com/linbaozhong/gentity/pkg/api/broker"
	"github.com/linbaozhong/gentity/pkg/api/iface"
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
	if id := ctx.GetHeader(broker.OperationID); len(id) == 0 {
		ctx.SetID(strconv.FormatInt(time.Now().UnixMilli(), 10))
	} else {
		ctx.SetID(id)
	}

	if ier, ok := arg.(iface.Initializer); ok {
		ier.Init()
	}
}
