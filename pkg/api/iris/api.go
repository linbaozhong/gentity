// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package api

import "github.com/linbaozhong/gentity/pkg/api/core"

// Initiate 初始化请求上下文
func Initiate(ctx Context, arg any) {
	core.Initiate(adapt(ctx), arg)
}

// Validate 校验参数
func Validate(arg any) error {
	return core.Validate(arg)
}
