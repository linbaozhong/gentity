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

package memcached

import (
	"context"
	"errors"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/conv"
	"strings"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type (
	// option 作用于非泛型配置，避免泛型实例的类型参数污染选项函数签名
	option func(o *cfg)

	cfg struct {
		prefix   string // key前缀
		duration int32  // 过期时间
	}

	memcached[V any] struct {
		mu sync.Mutex

		driver   *memcache.Client
		prefix   string // key前缀
		duration int32  // 过期时间
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
		o.duration = int32(duration.Seconds())
	}
}

// New creates an instance of Memcached cache driver
func New[V any](driver *memcache.Client, opts ...option) cachego.Cache[V] {
	var c cfg
	for _, opt := range opts {
		opt(&c)
	}

	obj := &memcached[V]{
		driver:   driver,
		prefix:   c.prefix,
		duration: c.duration,
	}
	return obj
}

// toBytes 将 V 序列化为 []byte，用于写入外部存储
func (m *memcached[V]) toBytes(value V) ([]byte, error) {
	return conv.Any2Bytes(value)
}

// fromBytes 将 []byte 反序列化为 V，用于从外部存储读取
func (m *memcached[V]) fromBytes(b []byte) (V, error) {
	var v V
	if len(b) == 0 {
		return v, nil
	}
	err := conv.Bytes2Any(b, &v)
	return v, err
}

// Contains checks if cached key exists in Memcached storage
func (m *memcached[V]) Contains(ctx context.Context, key string) bool {
	_, err := m.driver.Get(m.getKey(key))
	return err == nil
}

// ExistsOrSave 缓存不存在时，设置缓存，返回是否成功；缓存存在时，返回false
func (m *memcached[V]) ExistsOrSave(ctx context.Context, key string, value V, lifeTime ...time.Duration) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Contains(ctx, key) {
		return false
	}
	return m.Save(ctx, key, value, lifeTime...) == nil
}

// Delete the cached key from Memcached storage
func (m *memcached[V]) Delete(ctx context.Context, key string) error {
	return m.driver.Delete(m.getKey(key))
}

func (m *memcached[V]) PrefixDelete(ctx context.Context, prefix string) error {
	items, err := m.driver.GetMulti([]string{})
	if err != nil {
		return err
	}

	k := m.getKey(prefix)
	for _, item := range items {
		if matchPrefix(item.Key, k) {
			if err = m.driver.Delete(item.Key); err != nil {
				if !errors.Is(err, memcache.ErrCacheMiss) {
					return err
				}
			}
		}
	}
	return nil
}

// Fetch retrieves the cached value from key of the Memcached storage
func (m *memcached[V]) Fetch(ctx context.Context, key string) (V, error) {
	var zero V
	item, err := m.driver.Get(m.getKey(key))
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return zero, cachego.ErrCacheMiss
		}
		return zero, err
	}
	return m.fromBytes(item.Value)
}

// FetchMulti retrieves multiple cached value from keys of the Memcached storage
func (m *memcached[V]) FetchMulti(ctx context.Context, keys ...string) ([]V, error) {
	var ks []string

	if len(m.prefix) == 0 {
		ks = keys
	} else {
		ks = make([]string, 0, len(keys))
		for _, k := range keys {
			ks = append(ks, m.getKey(k))
		}
	}

	items, err := m.driver.GetMulti(ks)
	if err != nil {
		return nil, err
	}

	// 按请求顺序对齐结果，缺失的 key 填零值
	vals := make([]V, 0, len(ks))
	var zero V
	for _, k := range ks {
		item, ok := items[k]
		if !ok {
			vals = append(vals, zero)
			continue
		}
		v, err := m.fromBytes(item.Value)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}

	return vals, nil
}

// Flush removes all cached keys of the Memcached storage
func (m *memcached[V]) Flush(ctx context.Context) error {
	return m.driver.FlushAll()
}

// Save a value in Memcached storage by key
func (m *memcached[V]) Save(ctx context.Context, key string, value V, lifeTime ...time.Duration) error {
	val, err := m.toBytes(value)
	if err != nil {
		return err
	}

	duration := m.duration
	if len(lifeTime) > 0 {
		duration = int32(lifeTime[0].Seconds())
	}

	return m.driver.Set(
		&memcache.Item{
			Key:        m.getKey(key),
			Value:      val,
			Expiration: duration,
		},
	)
}

func (m *memcached[V]) getKey(key string) string {
	if m.prefix == "" {
		return key
	}
	return m.prefix + ":" + key
}

// matchPrefix 判断存储的 key 是否属于指定前缀命名空间。
// 边界为分隔符 ":"，因此 "user" 会命中 "user"、"user:1"，但不会误删 "username"。
func matchPrefix(storedKey, p string) bool {
	if storedKey == p {
		return true
	}
	return strings.HasPrefix(storedKey, p+":")
}
