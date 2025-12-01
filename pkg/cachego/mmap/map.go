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
	option func(o *syncMap)

	syncMapItem struct {
		data     []byte
		duration int64
	}

	syncMap struct {
		mu sync.Mutex

		storage  cmap.ConcurrentMap[string, any]
		prefix   string // key前缀
		duration int64  // 过期时间
	}
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *syncMap) {
		o.prefix = prefix
	}
}

// WithExpired 设置过期时间
func WithExpired(duration time.Duration) option {
	return func(o *syncMap) {
		o.duration = int64(duration.Seconds())
	}
}

// New creates an instance of SyncMap cache driver
func New(opts ...option) cachego.Cache {
	obj := &syncMap{storage: cmap.New[any]()}
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

func (sm *syncMap) read(ctx context.Context, key string) (syncMapItem, error) {
	var item syncMapItem
	v, ok := sm.storage.Get(sm.getKey(key))
	if !ok {
		return item, cachego.ErrCacheMiss
	}

	item = v.(syncMapItem)

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
func (sm *syncMap) Contains(ctx context.Context, key string) bool {
	return sm.storage.Has(sm.getKey(key))
}

// ExistsOrSave 缓存不存在时，设置缓存，返回是否成功；缓存存在时，返回false
func (sm *syncMap) ExistsOrSave(ctx context.Context, key string, value []byte, lifeTime ...time.Duration) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.Contains(ctx, key) {
		return false
	}
	return sm.Save(ctx, key, value, lifeTime...) == nil
}

// Delete the cached key from SyncMap storage
func (sm *syncMap) Delete(ctx context.Context, key string) error {
	sm.storage.Remove(sm.getKey(key))
	return nil
}

// PrefixDelete 按前缀删除
func (sm *syncMap) PrefixDelete(ctx context.Context, prefix string) error {
	k := sm.getKey(prefix)
	sm.storage.IterCb(func(key string, value any) {
		if strings.HasPrefix(key, k) {
			sm.storage.Remove(key)
		}
	})
	return nil
}

// Fetch retrieves the cached value from key of the SyncMap storage
func (sm *syncMap) Fetch(ctx context.Context, key string) ([]byte, error) {
	item, err := sm.read(ctx, key)
	if err != nil {
		return nil, err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the SyncMap storage
func (sm *syncMap) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
	vals := make([][]byte, 0, len(keys))
	for _, key := range keys {
		if b, err := sm.Fetch(ctx, key); err == nil {
			vals = append(vals, b)
		} else {
			vals = append(vals, nil)
		}
	}

	return vals, nil
}

// Flush removes all cached keys of the SyncMap storage
func (sm *syncMap) Flush(ctx context.Context) error {
	sm.storage.Clear()
	return nil
}

// Save a value in SyncMap storage by key
func (sm *syncMap) Save(ctx context.Context, key string, value []byte, lifeTime ...time.Duration) error {
	duration := time.Now().Unix() + sm.duration

	if len(lifeTime) > 0 {
		duration = time.Now().Unix() + int64(lifeTime[0].Seconds())
	}
	//buf, err := conv.Any2Bytes(value)
	//if err != nil {
	//	return err
	//}
	sm.storage.Set(sm.getKey(key), syncMapItem{value, duration})
	return nil
}

func (sm *syncMap) getKey(key string) string {
	if sm.prefix == "" {
		return key
	}
	return sm.prefix + ":" + key
}
