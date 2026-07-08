// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package api

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/api/core"
)

// GetResult 调用service处理get请求，并返回结果数据
func GetResult[A, B any](ctx Context,
	callService func(ctx context.Context, req *A, resp *B) error) (*B, error) {
	return core.GetResult(adapt(ctx),
		func(c context.Context, req *A, resp *B) error {
			return callService(c, req, resp)
		})
}

// PostResult 调用service处理post请求，并返回结果数据
func PostResult[A, B any](ctx Context,
	callService func(ctx context.Context, req *A, resp *B) error) (*B, error) {
	return core.PostResult(adapt(ctx),
		func(c context.Context, req *A, resp *B) error {
			return callService(c, req, resp)
		})
}

// StreamResult 调用service处理post请求，并返回结果数据
func StreamResult[A, B any](ctx Context,
	callService func(ctx context.Context, req *A, resp *B) error) (*B, error) {
	return core.StreamResult(adapt(ctx),
		func(c context.Context, req *A, resp *B) error {
			return callService(c, req, resp)
		})
}
