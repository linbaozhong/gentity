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

package sqlite

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/conv"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/util"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	// option 作用于非泛型配置，避免泛型实例的类型参数污染选项函数签名
	option func(o *cfg)

	cfg struct {
		interval time.Duration // 清理间隔
		prefix   string        // key前缀
		duration int64         // 过期时间
	}

	sqlite[V any] struct {
		db       *sql.DB
		lastTime time.Time     // 上次操作时间
		interval time.Duration // 空闲时长后清理
		name     string        // 缓存名称
		prefix   string        // key前缀
		duration int64         // 过期时间

		mu sync.Mutex
	}
)

const (
	cacheTableName       = "cache"
	cacheCleanupInterval = time.Minute
)

var (
	cacheDB   *sql.DB
	cacheOnce sync.Once
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *cfg) {
		o.prefix = prefix
	}
}

// WithInterval 设置清理间隔
func WithInterval(d time.Duration) option {
	return func(o *cfg) {
		o.interval = d
	}
}

// WithExpired 设置过期时间
func WithExpired(duration time.Duration) option {
	return func(o *cfg) {
		o.duration = int64(duration.Seconds())
	}
}

// New 创建一个sqlite缓存实例
func New[V any](ctx context.Context, name string, opts ...option) cachego.Cache[V] {
	cacheOnce.Do(func() {
		var err error
		err = os.MkdirAll("./cache", 0755)
		if err != nil {
			log.Fatal(err)
		}
		cacheDB, err = sql.Open("sqlite3", "file:cache/cache.db?cache=shared&mode=rwc&_journal_mode=WAL")
		if err != nil {
			log.Fatal(err)
		}
		cacheDB.Exec("PRAGMA synchronous = OFF")
	})

	var c cfg
	for _, opt := range opts {
		opt(&c)
	}

	obj := &sqlite[V]{
		db:       cacheDB,
		interval: cacheCleanupInterval,
		lastTime: time.Now(),
		prefix:   c.prefix,
		duration: c.duration,
	}
	if c.interval > 0 {
		obj.interval = c.interval
	}

	if name == "" {
		name = cacheTableName + "_" + util.GetRandLowerString(4)
	}
	obj.name = name

	obj.storage(ctx, obj.name)

	go obj.cleanup(ctx)

	return obj
}

// toBytes 将 V 序列化为 []byte，用于写入数据库
func (s *sqlite[V]) toBytes(value V) ([]byte, error) {
	return conv.Any2Bytes(value)
}

// fromBytes 将 []byte 反序列化为 V，用于从数据库读取
func (s *sqlite[V]) fromBytes(b []byte) (V, error) {
	var v V
	if len(b) == 0 {
		return v, nil
	}
	err := conv.Bytes2Any(b, &v)
	return v, err
}

func (s *sqlite[V]) Contains(ctx context.Context, key string) bool {
	var expires int64
	err := s.db.QueryRowContext(
		ctx,
		"SELECT expire FROM "+s.name+" WHERE key = ?",
		s.getKey(key),
	).Scan(&expires)
	if err != nil {
		return false
	}

	if expires == 0 || expires > time.Now().Unix() {
		return true
	}
	// 到期删除
	s.Delete(ctx, key)
	return false
}

// ExistsOrSave 缓存不存在时，设置缓存，返回是否成功；缓存存在时，返回false
func (s *sqlite[V]) ExistsOrSave(ctx context.Context, key string, value V, lifeTime ...time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Contains(ctx, key) {
		return false
	}

	buf, err := s.toBytes(value)
	if err != nil {
		return false
	}

	var dur int64
	if len(lifeTime) > 0 {
		if lifeTime[0] > 0 {
			dur = time.Now().Unix() + int64(lifeTime[0].Seconds())
		}
	} else if s.duration > 0 {
		dur = time.Now().Unix() + s.duration
	}
	// dur 为 0 表示永不过期

	result, err := s.db.ExecContext(ctx, "INSERT INTO "+s.name+"(value, expire,key) VALUES(?, ?, ?)", buf, dur, s.getKey(key))
	if err != nil {
		return false
	}
	n, _ := result.RowsAffected()
	return n > 0
}

func (s *sqlite[V]) Delete(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.name+" WHERE key = ?", s.getKey(key))
	return err
}

func (s *sqlite[V]) PrefixDelete(ctx context.Context, prefix string) error {
	k := s.getKey(prefix)
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.name+" WHERE key = ? OR key LIKE ?", k, k+":%")
	return err
}

func (s *sqlite[V]) Fetch(ctx context.Context, key string) (V, error) {
	var zero V
	var value []byte
	err := s.db.QueryRowContext(ctx, "SELECT value FROM "+s.name+" WHERE key = ? AND (expire = 0 OR expire > ?)",
		s.getKey(key),
		time.Now().Unix(),
	).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return zero, cachego.ErrCacheMiss
		}
		return zero, err
	}
	return s.fromBytes(value)
}

func (s *sqlite[V]) FetchMulti(ctx context.Context, keys ...string) ([]V, error) {
	l := len(keys)
	vals := make([]V, l)
	if l == 0 {
		return vals, nil
	}
	var (
		ks          = make([]any, 0, l)
		placeholder = strings.Repeat("?,", l)
	)
	// 为 in 查询做参数化
	for _, k := range keys {
		ks = append(ks, s.getKey(k))
	}
	ks = append(ks, time.Now().Unix())
	// 查询
	rows, err := s.db.QueryContext(ctx, "SELECT key, value FROM "+s.name+" WHERE key IN ("+
		placeholder[:len(placeholder)-1]+
		") AND (expire = 0 OR expire > ?)",
		ks...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 按请求顺序对齐结果，缺失的 key 保持零值
	byKey := make(map[string]V, l)
	for rows.Next() {
		var (
			k   string
			buf []byte
		)
		if rows.Scan(&k, &buf) == nil {
			if v, err := s.fromBytes(buf); err == nil {
				byKey[k] = v
			}
		}
	}
	for idx, k := range keys {
		if v, ok := byKey[s.getKey(k)]; ok {
			vals[idx] = v
		}
	}

	return vals, nil
}

// Flush 清空缓存
func (s *sqlite[V]) Flush(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.name)
	return err
}

func (s *sqlite[V]) Save(ctx context.Context, key string, value V, lifeTime ...time.Duration) error {
	buf, err := s.toBytes(value)
	if err != nil {
		return err
	}

	var (
		stmt *sql.Stmt
	)
	var dur int64
	if len(lifeTime) > 0 {
		if lifeTime[0] > 0 {
			dur = time.Now().Unix() + int64(lifeTime[0].Seconds())
		}
	} else if s.duration > 0 {
		dur = time.Now().Unix() + s.duration
	}
	// dur 为 0 表示永不过期
	// 查询是否存在
	if s.Contains(ctx, key) {
		stmt, err = s.db.PrepareContext(ctx, "UPDATE "+s.name+" SET value = ?, expire = ? WHERE key = ?")
	} else {
		stmt, err = s.db.PrepareContext(ctx, "INSERT INTO "+s.name+"(value, expire,key) VALUES(?, ?, ?)")
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, buf, dur, s.getKey(key))
	if err != nil {
		return err
	}
	return err
}

func (s *sqlite[V]) getKey(key string) string {
	if s.prefix == "" {
		return key
	}
	return s.prefix + ":" + key
}

func (s *sqlite[V]) storage(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS "`+name+`" (
			"key" TEXT NOT NULL DEFAULT '',
			"value" BLOB NOT NULL DEFAULT '',
			"expire" integer NOT NULL DEFAULT 0,
			PRIMARY KEY ("key")
		)`)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// cleanup 是一个定时运行的清理任务，用于删除过期对象。
func (s *sqlite[V]) cleanup(ctx context.Context) {
	// 创建定时器，用于定期清理过期对象。
	cleanTimer := time.NewTimer(s.interval)
	defer cleanTimer.Stop()

	for {
		select {
		case <-ctx.Done(): // 如果上下文被取消，退出并清理goroutine。
			log.Info("sqlite cache exit")
			return
		case <-cleanTimer.C:
			if time.Since(s.lastTime) > s.interval {
				s.db.ExecContext(ctx, "DELETE FROM "+s.name+" WHERE expire > 0 AND expire < ?", time.Now().Unix())
				s.lastTime = time.Now()
			}
			// 重置定时器。
			cleanTimer.Reset(s.interval)
		}
	}
}
