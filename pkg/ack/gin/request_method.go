// Copyright © 2023 SnowIM. All rights reserved.
// ...license same as before...

package ack

import (
	"context"

	"github.com/linbaozhong/gentity/pkg/ack/core"
)

// Get get请求：读取query。
func Get[A, B any](
	ctx Context,
	callService func(ctx context.Context, req *A, resp *B) error,
	after ...func(ctx Context, resp *B) error,
) error {
	return core.Get(adapt(ctx),
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

// Post post请求
func Post[A, B any](
	ctx Context,
	callService func(ctx context.Context, req *A, resp *B) error,
	after ...func(ctx Context, resp *B) error,
) error {
	return core.Post(adapt(ctx),
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

// Redirect 重定向
func Redirect[A any](ctx Context,
	callService func(ctx context.Context, req *A, resp *string) error,
) error {
	return core.Redirect(adapt(ctx),
		func(c context.Context, req *A, resp *string) error {
			return callService(c, req, resp)
		})
}

// Stream 流式请求
func Stream[A, B any](
	ctx Context,
	callService func(ctx context.Context, req *A, resp *B) error,
	after ...func(ctx Context, resp *B) error,
) error {
	return core.Stream(adapt(ctx),
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
