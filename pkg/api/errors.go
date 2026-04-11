package api

import (
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
	if ok := types.As(err, &appErr); ok && appErr != nil {
		switch types.ErrorType(appErr.Code) {
		case types.ErrorTypeDB:
			log.Errorf("【数据库错误】%s %s - IP: %s - 操作: %s - 错误: %v",
				method, path, ip, appErr.Op, appErr.Err)
		case types.ErrorTypeValidation:
			log.Warnf("【验证错误】%s %s - IP: %s - 操作: %s - 错误: %v",
				method, path, ip, appErr.Op, appErr.Err)
		case types.ErrorTypePermission:
			log.Warnf("【权限错误】%s %s - IP: %s - 操作: %s - 错误: %v",
				method, path, ip, appErr.Op, appErr.Err)
		case types.ErrorTypeParam:
			log.Errorf("【参数错误】%s %s - IP: %s - 操作: %s - 错误: %v",
				method, path, ip, appErr.Op, appErr.Err)
		case types.ErrorTypeNotFound:
			log.Infof("【未找到】%s %s - IP: %s - 操作: %s - 错误: %v",
				method, path, ip, appErr.Op, appErr.Err)
		case types.ErrorTypeUnknown:
			log.Panicf("【未知错误】%s %s - IP: %s - 操作: %s - 错误: %v",
				method, path, ip, appErr.Op, appErr.Err)
		default:
			log.Errorf("【未知错误】%s %s - IP: %s - 错误: %v",
				method, path, ip, err)
		}
	} else {
		// 不是 AppError，记录为普通错误
		log.Errorf("【系统错误】%s %s - IP: %s - 错误: %v",
			method, path, ip, err)
	}
}

func getErrorMessage(err error) string {
	// 根据错误类型记录不同级别的日志
	var appErr *types.Error
	if ok := types.As(err, &appErr); ok && appErr != nil {
		if appErr.Message != "" {
			return appErr.Message
		}
		switch types.ErrorType(appErr.Code) {
		case types.ErrorTypeDB:
			return "数据库错误"
		case types.ErrorTypeValidation:
			return "数据校验错误"
		case types.ErrorTypePermission:
			return "权限错误"
		case types.ErrorTypeParam:
			return "参数错误"
		case types.ErrorTypeNotFound:
			return "未找到"
		default:
			return "未知错误"
		}
	} else {
		return err.Error()
	}
}
