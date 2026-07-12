package core

import (
	"errors"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

func logError(c Context, err error) {
	path := c.Path()
	method := c.Method()
	ip := c.RemoteAddr()

	var appErr *types.Error
	if ok := errors.As(err, &appErr); ok && appErr != nil {
		log.Errorf("【%s】%s %s - IP: %s - 操作: %s - 错误: %v",
			appErr.Message, method, path, ip, appErr.Op, appErr.Err)
	} else {
		log.Errorf("【未知错误】%s %s - IP: %s - 错误: %v",
			method, path, ip, err)
	}
}
