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

package cachego

import (
	"context"
	"errors"
	"github.com/linbaozhong/gentity/pkg/conv"
	"github.com/linbaozhong/gentity/pkg/util"
	"time"
)

var (
	ErrCacheExpired = errors.New("cache expired")
	ErrCacheMiss    = errors.New("cache miss")
)

// Cache 泛型缓存接口
// V 表示缓存值类型，Fetch 返回 V，Save 接收 V。
// 内存型驱动(mmap)直接保存原始类型的 V，外部存储驱动(redis/memcached/sqlite)
// 内部通过 JSON 编解码进行序列化，因此 V 需可被 encoding/json 处理。
type Cache[V any] interface {
	Contains(ctx context.Context, key string) bool
	ExistsOrSave(ctx context.Context, key string, value V, lifeTime ...time.Duration) bool
	Delete(ctx context.Context, key string) error
	PrefixDelete(ctx context.Context, prefix string) error
	Fetch(ctx context.Context, key string) (V, error)
	FetchMulti(ctx context.Context, keys ...string) ([]V, error)
	Flush(ctx context.Context) error
	Save(ctx context.Context, key string, value V, lifeTime ...time.Duration) error
}

// GetHashKey 使用MemHash算法生成key
func GetHashKey(prefix string, key any) string {
	return prefix + util.HashString(conv.Any2String(key))
}
