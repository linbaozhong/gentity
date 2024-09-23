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

package redis

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/conv"
	"time"

	rd "github.com/redis/go-redis/v9"
)

type redis struct {
	driver rd.Cmdable
}

// New creates an instance of Redis cache driver
func New(driver rd.Cmdable) cachego.Cache {
	return &redis{driver}
}

// Contains checks if cached key exists in Redis storage
func (r *redis) Contains(ctx context.Context, key string) bool {
	i, _ := r.driver.Exists(ctx, key).Result()
	return i > 0
}

// Delete the cached key from Redis storage
func (r *redis) Delete(ctx context.Context, key string) error {
	return r.driver.Del(ctx, key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *redis) Fetch(ctx context.Context, key string) ([]byte, error) {
	return r.driver.Get(ctx, key).Bytes()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *redis) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
	items, err := r.driver.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	vals := make([][]byte, 0, len(items))
	for _, i := range items {
		val, err := conv.Interface2Bytes(i)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}

	return vals, nil
}

// Flush removes all cached keys of the Redis storage
func (r *redis) Flush(ctx context.Context) error {
	return r.driver.FlushAll(ctx).Err()
}

// Save a value in Redis storage by key
func (r *redis) Save(ctx context.Context, key string, value any, lifeTime time.Duration) error {
	return r.driver.Set(ctx, key, value, lifeTime).Err()
}
