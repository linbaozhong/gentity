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
	"fmt"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/log"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	option func(*sqlite)
	sqlite struct {
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

// WithName 设置缓存名称
func WithName(name string) option {
	return func(o *sqlite) {
		o.name = name
	}
}

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *sqlite) {
		o.prefix = prefix
	}
}

// WithInterval 设置清理间隔
func WithInterval(d time.Duration) option {
	return func(o *sqlite) {
		o.interval = d
	}
}

// WithExpired 设置过期时间
func WithExpired(duration time.Duration) option {
	return func(o *sqlite) {
		o.duration = time.Now().Unix() + int64(duration.Seconds())
	}
}

// New 创建一个sqlite缓存实例
func New(ctx context.Context, opts ...option) cachego.Cache {
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

	obj := &sqlite{
		db:       cacheDB,
		name:     cacheTableName,
		interval: cacheCleanupInterval,
	}

	for _, opt := range opts {
		opt(obj)
	}

	obj.storage(ctx, obj.name)

	go obj.cleanup(ctx)

	return obj
}

func (s *sqlite) Contains(ctx context.Context, key string) bool {
	var expires int64
	err := s.db.QueryRowContext(
		ctx,
		"SELECT expire FROM "+s.name+" WHERE key = ?",
		s.getKey(key),
	).Scan(&expires)
	if err != nil {
		return false
	}

	if expires > time.Now().Unix() {
		return true
	}
	// 到期删除
	s.Delete(ctx, key)
	return false
}

// ExistsOrSave 缓存不存在时，设置缓存，返回是否成功；缓存存在时，返回false
func (s *sqlite) ExistsOrSave(ctx context.Context, key string, value any, lifeTime ...time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Contains(ctx, key) {
		return false
	}

	duration := s.duration
	if len(lifeTime) > 0 {
		duration = time.Now().Unix() + int64(lifeTime[0].Seconds())
	}

	result, err := s.db.ExecContext(ctx, "INSERT INTO "+s.name+"(value, expire,key) VALUES(?, ?, ?)", value, time.Now().Unix()+duration, s.getKey(key))
	if err != nil {
		return false
	}
	n, _ := result.RowsAffected()
	return n > 0
}

func (s *sqlite) Delete(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.name+" WHERE key = ?", s.getKey(key))
	return err
}

func (s *sqlite) PrefixDelete(ctx context.Context, prefix string) error {
	k := s.getKey(prefix)
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.name+" WHERE key LIKE ?", k+"%")
	return err
}

func (s *sqlite) Fetch(ctx context.Context, key string) ([]byte, error) {
	s.lastTime = time.Now()

	var value []byte
	err := s.db.QueryRowContext(ctx, "SELECT value FROM "+s.name+" WHERE key = ? AND expire > ?",
		s.getKey(key),
		time.Now().Unix(),
	).Scan(&value)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cachego.ErrCacheMiss
		}
		return nil, err
	}
	return value, nil
}

func (s *sqlite) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
	s.lastTime = time.Now()
	var (
		l           = len(keys)
		vals        = make([][]byte, 0, l)
		ks          = make([]any, 0, l)
		placeholder = strings.Repeat("?,", l)
	)
	// 为 in 查询做参数化
	for _, k := range keys {
		ks = append(ks, s.getKey(k))
	}
	ks = append(ks, time.Now().Unix())
	// 查询
	rows, err := s.db.QueryContext(ctx, "SELECT value FROM "+s.name+" WHERE key IN ("+
		placeholder[:len(placeholder)-1]+
		") AND expire > ?",
		ks...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 扫描结果
	for rows.Next() {
		var value []byte
		if rows.Scan(&value) == nil {
			vals = append(vals, value)
		}
	}

	return vals, nil
}

// Flush 清空缓存
func (s *sqlite) Flush(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+s.name)
	return err
}

func (s *sqlite) Save(ctx context.Context, key string, value any, lifeTime ...time.Duration) error {
	var (
		stmt *sql.Stmt
		err  error
	)
	duration := s.duration
	if len(lifeTime) > 0 {
		duration = time.Now().Unix() + int64(lifeTime[0].Seconds())
	}
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

	_, err = stmt.ExecContext(ctx, value, time.Now().Unix()+duration, s.getKey(key))
	if err != nil {
		return err
	}
	return err
}

func (s *sqlite) getKey(key string) string {
	if s.prefix == "" {
		return key
	}
	return s.prefix + ":" + key
}

func (s *sqlite) storage(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS "`+name+`" (
			"key" TEXT NOT NULL DEFAULT '',
			"value" TEXT NOT NULL DEFAULT '',
			"expire" integer NOT NULL DEFAULT 0,
			PRIMARY KEY ("key")
		)`)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// cleanup 是一个定时运行的清理任务，用于删除过期对象。
func (s *sqlite) cleanup(ctx context.Context) {
	// 创建定时器，用于定期清理过期对象。
	cleanTimer := time.NewTimer(s.interval)
	defer cleanTimer.Stop()

	for {
		select {
		case <-ctx.Done(): // 如果上下文被取消，退出并清理goroutine。
			fmt.Println("cleanup exit")
			return
		case <-cleanTimer.C:
			if time.Since(s.lastTime) > s.interval {
				s.db.ExecContext(ctx, "DELETE FROM "+s.name+" WHERE expire < ?", time.Now().Unix())
			}
			// 重置定时器。
			cleanTimer.Reset(s.interval)
		}
	}
}
