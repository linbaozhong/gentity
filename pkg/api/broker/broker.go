package broker

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/api/iface"
)

const (
	OperationID = "operation_id"
)

var (
	// 全局context，用于支持外部调用。最好在程序启动时引用，否则可能造成panic。
	// 在程序退出时，需要调用Cancel()。
	Context, Cancel = context.WithCancel(context.Background())

	ServiceCloses = make([]any, 0)
)

func Validate(arg any) error {
	if checker, ok := arg.(iface.Checker); ok {
		return checker.Check()
	}
	return nil
}

// RegisterServiceCloser 注册服务关闭器
func RegisterServiceCloser(s iface.IServiceCloser) {
	ServiceCloses = append(ServiceCloses, s)
}

// Close 关闭所有服务
func Close() {
	Cancel()
	for _, s := range ServiceCloses {
		if s, ok := s.(iface.IServiceCloser); ok {
			s.Close()
		}
	}
}
