package lib

import (
	ack "github.com/linbaozhong/gentity/pkg/ack/iris"
	"one/internal/constant"
)

func AuthMiddleware() ack.Handler {
	return func(c ack.Context) {
		token := c.GetHeader(constant.Authorization)
		if token == "" {
			ack.Fail(c, constant.ErrAuthorizationNotFound)
			return
		}
		c.Next()
	}
}
