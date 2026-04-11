package api

import (
	"errors"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

// logError 记录错误日志
func logError(c Context, err error) {
	// 获取请求信息
	path := c.Path()
	method := c.Method()
	ip := c.RemoteAddr()

	// 根据错误类型记录不同级别的日志
	var appErr *types.Error
	if ok := errors.As(err, &appErr); ok && appErr != nil {
		log.Errorf("【%s】%s %s - IP: %s - 操作: %s - 错误: %v",
			appErr.Message, method, path, ip, appErr.Op, appErr.Err)
	} else {
		// 不是 AppError，记录为普通错误
		log.Errorf("【未知错误】%s %s - IP: %s - 错误: %v",
			method, path, ip, err)
	}
}
