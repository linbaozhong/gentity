// Copyright © 2023 SnowIM. All rights reserved.
// ...license same as before...

package ack

import (
	"github.com/linbaozhong/gentity/pkg/ack/core"
)

// Fail 写入错误响应
func Fail(c Context, e error, args ...any) error {
	return core.Fail(adapt(c), e, args...)
}

// Ok 写入成功响应，支持 GET 缓存
func Ok(c Context, args ...any) error {
	return core.Ok(adapt(c), args...)
}

// SendLocalFile 发送本地文件
func SendLocalFile(c Context, path, name string) error {
	return core.SendLocalFile(adapt(c), path, name)
}

// SendUrlFile 发送url文件
func SendUrlFile(c Context, url, name string) error {
	return core.SendUrlFile(adapt(c), url, name)
}
