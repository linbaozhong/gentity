// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package ack

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/ack/internal/core"
	"github.com/linbaozhong/gentity/pkg/cachego"
)

// DefaultIdempotencyConfig 返回默认配置
func DefaultIdempotencyConfig(cache cachego.Cache) *core.IdempotencyConfig {
	return core.DefaultIdempotencyConfig(cache)
}

// PostIdempotent 支持幂等的 post 请求
func PostIdempotent[A, B any](
	ctx Context,
	config *core.IdempotencyConfig,
	callService func(ctx context.Context, req *A, resp *B) error,
	after ...func(ctx Context, resp *B) error,
) error {
	return core.PostIdempotent(adapt(ctx), config,
		func(c context.Context, req *A, resp *B) error {
			return callService(c, req, resp)
		},
		func(c core.Context, resp *B) error {
			if len(after) > 0 {
				return after[0](ctx, resp)
			}
			return nil
		})
}
