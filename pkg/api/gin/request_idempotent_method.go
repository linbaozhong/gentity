// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package api

import (
	"github.com/linbaozhong/gentity/pkg/api/core"
	"github.com/linbaozhong/gentity/pkg/cachego"
)

// IdempotencyConfig 幂等性配置
type IdempotencyConfig = core.IdempotencyConfig

// DefaultIdempotencyConfig 返回默认配置
func DefaultIdempotencyConfig(cache cachego.Cache) *IdempotencyConfig {
	return core.DefaultIdempotencyConfig(cache)
}
