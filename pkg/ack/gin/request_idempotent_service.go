// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package ack

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/ack/internal/core"
)

// PostIdempotentResult 支持幂等的 post 请求，返回结果数据
func PostIdempotentResult[A, B any](
	ctx Context,
	config *core.IdempotencyConfig,
	callService func(ctx context.Context, req *A, resp *B) error,
) (*B, error) {
	return core.PostIdempotentResult(adapt(ctx), config,
		func(c context.Context, req *A, resp *B) error {
			return callService(c, req, resp)
		})
}
