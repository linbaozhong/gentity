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

package mmap

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/cachego"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strings"
	"sync"
	"time"
)

type (
	// option 作用于非泛型配置，避免泛型实例的类型参数污染选项函数签名
	option func(o *cfg)

	cfg struct {
		prefix   string // key前缀
		duration int64  // 过期时间
	}

	syncMapItem[V any] struct {
		data     V
		duration int64
	}

	syncMap[V any] struct {
		mu sync.Mutex

		storage  cmap.ConcurrentMap[string, syncMapItem[V]]
		prefix   string // key前缀
		duration int64  // 过期时间
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
		o.duration = int64(duration.Seconds())
	}
}

// New creates an instance of SyncMap cache driver
func New[V any](opts ...option) cachego.Cache[V] {
	var c cfg
	for _, opt := range opts {
		opt(&c)
	}

	obj := &syncMap[V]{
		prefix:   c.prefix,
		duration: c.duration,
		storage:  cmap.New[syncMapItem[V]](),
	}
	return obj
}

func (sm *syncMap[V]) read(ctx context.Context, key string) (syncMapItem[V], error) {
	var item syncMapItem[V]
	v, ok := sm.storage.Get(sm.getKey(key))
	if !ok {
		return item, cachego.ErrCacheMiss
	}

	item = v

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		_ = sm.Delete(ctx, key)
		return item, cachego.ErrCacheMiss
	}

	return item, nil
}

// Contains checks if cached key exists in SyncMap storage
func (sm *syncMap[V]) Contains(ctx context.Context, key string) bool {
	return sm.storage.Has(sm.getKey(key))
}

// ExistsOrSave 缓存不存在时，设置缓存，返回是否成功；缓存存在时，返回false
func (sm *syncMap[V]) ExistsOrSave(ctx context.Context, key string, value V, lifeTime ...time.Duration) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.Contains(ctx, key) {
		return false
	}
	return sm.Save(ctx, key, value, lifeTime...) == nil
}

// Delete the cached key from SyncMap storage
func (sm *syncMap[V]) Delete(ctx context.Context, key string) error {
	sm.storage.Remove(sm.getKey(key))
	return nil
}

// PrefixDelete 按前缀删除
func (sm *syncMap[V]) PrefixDelete(ctx context.Context, prefix string) error {
	k := sm.getKey(prefix)
	// 注意：cmap 的 IterCb 在回调中持有 shard 读锁，不能在回调内调用 Remove（写锁）否则死锁。
	// 因此先收集待删 key，迭代结束后再删除。
	var toRemove []string
	sm.storage.IterCb(func(key string, value syncMapItem[V]) {
		if matchPrefix(key, k) {
			toRemove = append(toRemove, key)
		}
	})
	for _, key := range toRemove {
		sm.storage.Remove(key)
	}
	return nil
}

// Fetch retrieves the cached value from key of the SyncMap storage
func (sm *syncMap[V]) Fetch(ctx context.Context, key string) (V, error) {
	item, err := sm.read(ctx, key)
	if err != nil {
		var zero V
		return zero, err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the SyncMap storage
func (sm *syncMap[V]) FetchMulti(ctx context.Context, keys ...string) ([]V, error) {
	vals := make([]V, 0, len(keys))
	for _, key := range keys {
		if b, err := sm.Fetch(ctx, key); err == nil {
			vals = append(vals, b)
		} else {
			var zero V
			vals = append(vals, zero)
		}
	}

	return vals, nil
}

// Flush removes all cached keys of the SyncMap storage
func (sm *syncMap[V]) Flush(ctx context.Context) error {
	sm.storage.Clear()
	return nil
}

// Save a value in SyncMap storage by key
func (sm *syncMap[V]) Save(ctx context.Context, key string, value V, lifeTime ...time.Duration) error {
	var duration int64
	if len(lifeTime) > 0 {
		if lifeTime[0] > 0 {
			duration = time.Now().Unix() + int64(lifeTime[0].Seconds())
		}
	} else if sm.duration > 0 {
		duration = time.Now().Unix() + sm.duration
	}
	// duration 为 0 表示永不过期（read 中当 item.duration == 0 直接返回）
	sm.storage.Set(sm.getKey(key), syncMapItem[V]{data: value, duration: duration})
	return nil
}

func (sm *syncMap[V]) getKey(key string) string {
	if sm.prefix == "" {
		return key
	}
	return sm.prefix + ":" + key
}

// matchPrefix 判断存储的 key 是否属于指定前缀命名空间。
// 边界为分隔符 ":"，因此 "user" 会命中 "user"、"user:1"，但不会误删 "username"。
func matchPrefix(storedKey, p string) bool {
	if storedKey == p {
		return true
	}
	return strings.HasPrefix(storedKey, p+":")
}
