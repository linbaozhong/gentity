// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package ack

import (
	"time"

	"github.com/linbaozhong/gentity/pkg/ack/core"
)

// ReadCache 读取缓存，如果命中则直接写入响应并返回 true。
func ReadCache(ctx Context, lefetime ...time.Duration) bool {
	return core.ReadCache(adapt(ctx), lefetime...)
}
