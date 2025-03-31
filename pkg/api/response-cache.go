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
	"github.com/linbaozhong/gentity/pkg/cachego/mmap"
	"net/http"
	"time"
)

const (
	hasCacheKey      = "_HAS_CACHE_"
	authorizationKey = "Authorization"
)

type cacheKey struct {
	Key      string
	Duration time.Duration
}

var (
	respCache = mmap.New(mmap.WithExpired(time.Second * 30))
)

// ReadCache 读取缓存。
// 注意：如果缓存存在，会直接返回，不会再执行后续的逻辑。
func ReadCache(ctx Context, lefetime ...time.Duration) bool {
	_url := ctx.Request().URL
	_vals := _url.Query()
	_vals.Set("_t", ctx.GetHeader(authorizationKey))
	_key := _url.Path + "?" + _vals.Encode()

	duration := time.Second * 30
	if len(lefetime) > 0 {
		duration = lefetime[0]
	}
	ctx.Values().Set(hasCacheKey, cacheKey{
		Key:      _key,
		Duration: duration,
	})

	buf, e := respCache.Fetch(ctx, _key)
	if e != nil {
		return false
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.StatusCode(http.StatusOK)
	if _, e = ctx.Write(buf); e != nil {
		return false
	}
	ctx.EndRequest()
	return true
}

func setCache(ctx Context, resp any) {
	if _cacheKey := ctx.Values().Get(hasCacheKey); _cacheKey != nil {
		if _key, ok := _cacheKey.(cacheKey); ok {
			respCache.Save(ctx, _key.Key, resp, _key.Duration)
		}
	}
}
