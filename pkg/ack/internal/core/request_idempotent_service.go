// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package core

import (
	"context"
	"encoding/json"
	"fmt"
)

// PostIdempotent 支持幂等的 post 请求
func PostIdempotent[A, B any](
	ctx Context,
	config *IdempotencyConfig,
	callService func(ctx context.Context, req *A, resp *B) error,
	after ...func(ctx Context, resp *B) error,
) error {
	resultPtr, err := PostIdempotentResult(ctx, config, callService)
	if err != nil {
		return Fail(ctx, err)
	}
	if len(after) > 0 && resultPtr != nil {
		after[0](ctx, resultPtr)
	}
	return Ok(ctx, resultPtr)
}

// PostIdempotentResult 支持幂等的 post 请求，返回结果数据
func PostIdempotentResult[A, B any](
	ctx Context,
	config *IdempotencyConfig,
	callService func(ctx context.Context, req *A, resp *B) error,
) (*B, error) {
	var (
		req  A
		resp B
	)

	keyFn := config.KeyFunc
	if keyFn == nil {
		keyFn = defaultIdempotencyKey
	}
	cacheKey := "idemp:" + keyFn(ctx)

	// 快速路径
	if data, err := config.Cache.Fetch(ctx, cacheKey); err == nil {
		if json.Unmarshal(data, &resp) == nil {
			return &resp, nil
		}
	}

	// 单飞保护
	val, err, _ := idempotencyGroup.Do(cacheKey, func() (interface{}, error) {
		_, e := serviceContext(ctx, &req, &resp, readPostRequest[A], callService)
		if e != nil {
			return nil, e
		}
		if data, marshalErr := json.Marshal(&resp); marshalErr == nil {
			config.Cache.Save(ctx, cacheKey, data, config.ExpireIn)
		}
		return &resp, nil
	})

	if err != nil {
		return nil, err
	}

	if b, ok := val.(*B); ok {
		return b, nil
	}
	return nil, fmt.Errorf("idempotency: unexpected type %T", val)
}
