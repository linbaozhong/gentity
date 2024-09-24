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
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type (
	option func(o *memcached)

	memcached struct {
		driver *memcache.Client
		prefix string // key前缀
	}
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *memcached) {
		o.prefix = prefix
	}
}

// New creates an instance of Memcached cache driver
func New(driver *memcache.Client, opts ...option) cachego.Cache {
	obj := &memcached{driver: driver}
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

// Contains checks if cached key exists in Memcached storage
func (m *memcached) Contains(ctx context.Context, key string) bool {
	_, err := m.Fetch(ctx, m.getKey(key))
	return err == nil
}

// Delete the cached key from Memcached storage
func (m *memcached) Delete(ctx context.Context, key string) error {
	return m.driver.Delete(m.getKey(key))
}

func (m *memcached) PrefixDelete(ctx context.Context, prefix string) error {
	items, err := m.driver.GetMulti([]string{})
	if err != nil {
		return err
	}

	k := m.getKey(prefix)
	for _, item := range items {
		if strings.HasPrefix(item.Key, k) {
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
func (m *memcached) Fetch(ctx context.Context, key string) ([]byte, error) {
	item, err := m.driver.Get(m.getKey(key))
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, cachego.ErrCacheMiss
		}
		return nil, err
	}
	return item.Value, nil
}

// FetchMulti retrieves multiple cached value from keys of the Memcached storage
func (m *memcached) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
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

	vals := make([][]byte, 0, len(items))
	for _, i := range items {
		vals = append(vals, i.Value)
	}

	return vals, nil
}

// Flush removes all cached keys of the Memcached storage
func (m *memcached) Flush(ctx context.Context) error {
	return m.driver.FlushAll()
}

// Save a value in Memcached storage by key
func (m *memcached) Save(ctx context.Context, key string, value any, lifeTime time.Duration) error {
	val, err := conv.Interface2Bytes(value)
	if err != nil {
		return err
	}
	return m.driver.Set(
		&memcache.Item{
			Key:        m.getKey(key),
			Value:      val,
			Expiration: int32(lifeTime.Seconds()),
		},
	)
}

func (m *memcached) getKey(key string) string {
	if m.prefix == "" {
		return key
	}
	return m.prefix + ":" + key
}