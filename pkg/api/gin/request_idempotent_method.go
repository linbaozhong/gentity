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
	// 参数：原始请求 body []byte
	// 返回：幂等键值
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
	body, _ := io.ReadAll(ctx.Request.Body)
	return util.HashString(ctx.Request.URL.Path + string(body))
}

// PostIdempotent 支持幂等的 post 请求
// 自动基于请求 body 生成幂等键，无需客户端传额外参数
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
