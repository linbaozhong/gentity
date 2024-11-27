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
	"strconv"
	"time"
)

var (
	ErrCacheExpired = errors.New("cache expired")
	ErrCacheMiss    = errors.New("cache miss")
)

type Cache interface {
	Contains(ctx context.Context, key string) bool
	ContainsOrSave(ctx context.Context, key string, value any, lifeTime time.Duration) bool
	Delete(ctx context.Context, key string) error
	PrefixDelete(ctx context.Context, prefix string) error
	Fetch(ctx context.Context, key string) ([]byte, error)
	FetchMulti(ctx context.Context, keys ...string) ([][]byte, error)
	Flush(ctx context.Context) error
	Save(ctx context.Context, key string, value any, lifeTime time.Duration) error
}

// Hash 使用MemHash算法
func Hash(key any) string {
	return strconv.FormatUint(util.MemHashString(conv.Any2String(key)), 10)
}

// GetHashKey 使用MemHash算法生成key
func GetHashKey(prefix string, key any) string {
	return prefix + Hash(key)
}
