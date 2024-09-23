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
	"go/types"
	"time"
)

var (
	ErrCacheExpired = errors.New("cache expired")
)

type CacheValueType interface {
	int | float64 | string | bool | []byte | types.Map | types.Slice | types.Struct
}

type Cache interface {
	Contains(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
	Fetch(ctx context.Context, key string) ([]byte, error)
	FetchMulti(ctx context.Context, keys ...string) ([][]byte, error)
	Flush(ctx context.Context) error
	Save(ctx context.Context, key string, value any, lifeTime time.Duration) error
}
