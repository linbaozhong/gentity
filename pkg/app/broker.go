package app

import (
	"context"
)

type IServiceLauncher interface {
	Launch() error
}
type IServiceCloser interface {
	Close() error
}

var (
	// 全局context，用于支持外部调用。最好在程序启动时引用，否则可能造成panic。
	// 在程序退出时，需要调用Cancel()。
	Context, Cancel = context.WithCancel(context.Background())

	serviceCloses   = make([]any, 0)
	serviceLaunches = make([]any, 0)
)

// RegisterServiceLauncher 注册服务启动器
func RegisterServiceLauncher(s IServiceLauncher) {
	serviceLaunches = append(serviceLaunches, s)
}

// Launch 启动所有服务
func Launch() {
	for _, s := range serviceLaunches {
		if s, ok := s.(IServiceLauncher); ok {
			s.Launch()
		}
	}
}

// RegisterServiceCloser 注册服务关闭器
func RegisterServiceCloser(s IServiceCloser) {
	serviceCloses = append(serviceCloses, s)
}

// Close 关闭所有服务
func Close() {
	Cancel()
	for _, s := range serviceCloses {
		if s, ok := s.(IServiceCloser); ok {
			s.Close()
		}
	}
}
