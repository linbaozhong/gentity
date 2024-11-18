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
	"errors"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/conv"
	"time"

	rd "github.com/redis/go-redis/v9"
)

type (
	option func(o *redis)

	redis struct {
		driver rd.Cmdable
		prefix string // key前缀
	}
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *redis) {
		o.prefix = prefix
	}
}

// New creates an instance of Redis cache driver
func New(driver rd.Cmdable, opts ...option) cachego.Cache {
	obj := &redis{driver: driver}
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

// Contains checks if cached key exists in Redis storage
func (r *redis) Contains(ctx context.Context, key string) bool {
	i, _ := r.driver.Exists(ctx, r.getKey(key)).Result()
	return i > 0
}

// Delete the cached key from Redis storage
func (r *redis) Delete(ctx context.Context, key string) error {
	return r.driver.Del(ctx, r.getKey(key)).Err()
}

// PrefixDelete 按前缀删除
func (r *redis) PrefixDelete(ctx context.Context, prefix string) error {
	iter := r.driver.Scan(ctx, 0, r.getKey(prefix), 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		// 使用 DEL 命令删除 key
		if err := r.driver.Del(ctx, key).Err(); err != nil {
			if err != rd.Nil {
				return err
			}
		}
	}
	return iter.Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *redis) Fetch(ctx context.Context, key string) ([]byte, error) {
	b, e := r.driver.Get(ctx, r.getKey(key)).Bytes()
	if e != nil {
		if errors.Is(e, rd.Nil) {
			return nil, cachego.ErrCacheMiss
		}
		return nil, e
	}

	return b, nil
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *redis) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
	var ks []string

	if len(r.prefix) == 0 {
		ks = keys
	} else {
		ks = make([]string, 0, len(keys))
		for _, k := range keys {
			ks = append(ks, r.getKey(k))
		}
	}

	items, err := r.driver.MGet(ctx, ks...).Result()
	if err != nil {
		return nil, err
	}

	vals := make([][]byte, 0, len(items))
	for _, i := range items {
		val, err := conv.Any2Bytes(i)
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
	return r.driver.Set(ctx, r.getKey(key), value, lifeTime).Err()
}

func (r *redis) getKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return r.prefix + ":" + key
}
