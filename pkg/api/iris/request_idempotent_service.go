// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package api

import (
	"context"

	"github.com/linbaozhong/gentity/pkg/api/core"
)

// PostIdempotent 支持幂等的 post 请求
func PostIdempotent[A, B any](
	ctx Context,
	config *IdempotencyConfig,
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

// PostIdempotentResult 支持幂等的 post 请求，返回结果数据
func PostIdempotentResult[A, B any](
	ctx Context,
	config *IdempotencyConfig,
	callService func(ctx context.Context, req *A, resp *B) error,
) (*B, error) {
	return core.PostIdempotentResult(adapt(ctx), config,
		func(c context.Context, req *A, resp *B) error {
			return callService(c, req, resp)
		})
}
