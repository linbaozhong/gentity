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

package sync

import (
	"context"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/conv"
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
		storage *sync.Map
		prefix  string // key前缀
	}
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *syncMap) {
		o.prefix = prefix
	}
}

// New creates an instance of SyncMap cache driver
func New(opts ...option) cachego.Cache {
	obj := &syncMap{storage: &sync.Map{}}
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

func (sm *syncMap) read(ctx context.Context, key string) (*syncMapItem, error) {
	v, ok := sm.storage.Load(key)
	if !ok {
		return nil, cachego.ErrCacheMiss
	}

	item := v.(*syncMapItem)

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		_ = sm.Delete(ctx, key)
		return nil, cachego.ErrCacheMiss
	}

	return item, nil
}

// Contains checks if cached key exists in SyncMap storage
func (sm *syncMap) Contains(ctx context.Context, key string) bool {
	_, ok := sm.storage.Load(sm.getKey(key))
	return ok
}

// Delete the cached key from SyncMap storage
func (sm *syncMap) Delete(ctx context.Context, key string) error {
	sm.storage.Delete(sm.getKey(key))
	return nil
}

// PrefixDelete 按前缀删除
func (sm *syncMap) PrefixDelete(ctx context.Context, prefix string) error {
	k := sm.getKey(prefix)
	sm.storage.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), k) {
			sm.storage.Delete(key)
		}
		return true
	})
	return nil
}

// Fetch retrieves the cached value from key of the SyncMap storage
func (sm *syncMap) Fetch(ctx context.Context, key string) ([]byte, error) {
	item, err := sm.read(ctx, sm.getKey(key))
	if err != nil {
		return nil, err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the SyncMap storage
func (sm *syncMap) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
	vals := make([][]byte, 0, len(keys))
	for _, key := range keys {
		if b, err := sm.Fetch(ctx, sm.getKey(key)); err == nil {
			vals = append(vals, b)
		} else {
			vals = append(vals, nil)
		}
	}

	return vals, nil
}

// Flush removes all cached keys of the SyncMap storage
func (sm *syncMap) Flush(ctx context.Context) error {
	sm.storage = &sync.Map{}
	return nil
}

// Save a value in SyncMap storage by key
func (sm *syncMap) Save(ctx context.Context, key string, value any, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}
	b, err := conv.Interface2Bytes(value)
	if err != nil {
		return err
	}
	sm.storage.Store(sm.getKey(key), &syncMapItem{b, duration})
	return nil
}

func (sm *syncMap) getKey(key string) string {
	if sm.prefix == "" {
		return key
	}
	return sm.prefix + ":" + key
}
