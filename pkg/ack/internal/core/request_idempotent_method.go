// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package core

import (
	"io"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/util"
)

// idempotencyGroup 包级单飞组
var idempotencyGroup = &singleflight.Group{}

// IdempotencyConfig 幂等性配置
type IdempotencyConfig struct {
	// Cache 存储后端（必填）
	Cache cachego.Cache

	// ExpireIn 缓存过期时间，默认 24 小时
	ExpireIn time.Duration

	// KeyFunc 自定义幂等键生成函数，默认为 Hash(请求body)
	KeyFunc func(Context) string
}

// DefaultIdempotencyConfig 返回默认配置
func DefaultIdempotencyConfig(cache cachego.Cache) *IdempotencyConfig {
	return &IdempotencyConfig{
		Cache:    cache,
		ExpireIn: 24 * time.Hour,
		KeyFunc:  defaultIdempotencyKey,
	}
}

// defaultIdempotencyKey 默认幂等键：对请求 body 做 MemHash
func defaultIdempotencyKey(ctx Context) string {
	body, _ := io.ReadAll(ctx.Request().Body)
	return util.HashString(ctx.Request().URL.Path + string(body))
}
