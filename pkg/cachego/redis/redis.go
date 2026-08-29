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
	// option 作用于非泛型配置，避免泛型实例的类型参数污染选项函数签名
	option func(o *cfg)

	cfg struct {
		prefix   string        // key前缀
		duration time.Duration // 过期时间
	}

	redis[V any] struct {
		driver   rd.Cmdable
		prefix   string        // key前缀
		duration time.Duration // 过期时间
	}
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *cfg) {
		o.prefix = prefix
	}
}

// WithExpired 设置过期时间
func WithExpired(duration time.Duration) option {
	return func(o *cfg) {
		o.duration = duration
	}
}

// New creates an instance of Redis cache driver
func New[V any](driver rd.Cmdable, opts ...option) cachego.Cache[V] {
	var c cfg
	for _, opt := range opts {
		opt(&c)
	}

	obj := &redis[V]{
		driver:   driver,
		prefix:   c.prefix,
		duration: c.duration,
	}
	return obj
}

// toBytes 将 V 序列化为 []byte，用于写入外部存储
func (r *redis[V]) toBytes(value V) ([]byte, error) {
	return conv.Any2Bytes(value)
}

// fromBytes 将 []byte 反序列化为 V，用于从外部存储读取
func (r *redis[V]) fromBytes(b []byte) (V, error) {
	var v V
	if len(b) == 0 {
		return v, nil
	}
	err := conv.Bytes2Any(b, &v)
	return v, err
}

// Contains checks if cached key exists in Redis storage
func (r *redis[V]) Contains(ctx context.Context, key string) bool {
	i, _ := r.driver.Exists(ctx, r.getKey(key)).Result()
	return i > 0
}

// ExistsOrSave 缓存不存在时，设置缓存，返回是否成功；缓存存在时，返回false
func (r *redis[V]) ExistsOrSave(ctx context.Context, key string, value V, lifeTime ...time.Duration) bool {
	duration := r.duration
	if len(lifeTime) > 0 {
		duration = lifeTime[0]
	}
	buf, err := r.toBytes(value)
	if err != nil {
		return false
	}
	return r.driver.SetNX(ctx, r.getKey(key), buf, duration).Val()
}

// Delete the cached key from Redis storage
func (r *redis[V]) Delete(ctx context.Context, key string) error {
	return r.driver.Del(ctx, r.getKey(key)).Err()
}

// PrefixDelete 按前缀删除
func (r *redis[V]) PrefixDelete(ctx context.Context, prefix string) error {
	iter := r.driver.Scan(ctx, 0, r.getKey(prefix)+":*", 0).Iterator()
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
func (r *redis[V]) Fetch(ctx context.Context, key string) (V, error) {
	var zero V
	b, e := r.driver.Get(ctx, r.getKey(key)).Bytes()
	if e != nil {
		if errors.Is(e, rd.Nil) {
			return zero, cachego.ErrCacheMiss
		}
		return zero, e
	}

	return r.fromBytes(b)
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *redis[V]) FetchMulti(ctx context.Context, keys ...string) ([]V, error) {
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

	// MGet 按请求顺序返回，缺失的 key 对应 nil，这里填充零值
	vals := make([]V, 0, len(items))
	var zero V
	for _, i := range items {
		if i == nil {
			vals = append(vals, zero)
			continue
		}
		b, err := conv.Any2Bytes(i)
		if err != nil {
			return nil, err
		}
		v, err := r.fromBytes(b)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}

	return vals, nil
}

// Flush removes all cached keys of the Redis storage
func (r *redis[V]) Flush(ctx context.Context) error {
	return r.driver.FlushAll(ctx).Err()
}

// Save a value in Redis storage by key
func (r *redis[V]) Save(ctx context.Context, key string, value V, lifeTime ...time.Duration) error {
	duration := r.duration
	if len(lifeTime) > 0 {
		duration = lifeTime[0]
	}
	buf, err := r.toBytes(value)
	if err != nil {
		return err
	}
	return r.driver.Set(ctx, r.getKey(key), buf, duration).Err()
}

func (r *redis[V]) getKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return r.prefix + ":" + key
}
