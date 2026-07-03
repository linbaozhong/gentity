// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"encoding/json"
)

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

	return val.(*B), nil
}
