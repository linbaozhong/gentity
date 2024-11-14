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
	"github.com/linbaozhong/gentity/pkg/log"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
	"time"
)

type (
	option func(*sqlite)
	sqlite struct {
		db     *sql.DB
		prefix string // key前缀
	}
)

const (
	cacheTableName = "cache"
)

var (
	cacheDB   *sql.DB
	cacheOnce sync.Once
)

// WithPrefix 设置key前缀
func WithPrefix(prefix string) option {
	return func(o *sqlite) {
		o.prefix = prefix
	}
}

// New 创建一个sqlite缓存实例
func New(opts ...option) cachego.Cache {
	cacheOnce.Do(func() {
		var err error
		cacheDB, err = sql.Open("sqlite3", "file:cache.db?cache=shared&mode=rwc&_journal_mode=WAL")
		if err != nil {
			log.Fatal(err)
		}
		_, err = cacheDB.Exec(`CREATE TABLE IF NOT EXISTS "cache" (
			"key" TEXT NOT NULL DEFAULT '',
			"value" TEXT NOT NULL DEFAULT '',
			"expire" integer NOT NULL DEFAULT 0,
			PRIMARY KEY ("key")
		)`)

		if err != nil {
			log.Fatal(err)
		}
	})

	obj := &sqlite{db: cacheDB}
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

func (s *sqlite) Contains(ctx context.Context, key string) bool {
	var expires int64
	err := s.db.QueryRowContext(ctx, "SELECT expire FROM "+cacheTableName+" WHERE key = ?",
		s.getKey(key)).Scan(&expires)
	if err != nil {
		return false
	}

	if expires < time.Now().Unix() && expires != 0 {
		s.Delete(ctx, key)
		return false
	}
	return true
}

func (s *sqlite) Delete(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+cacheTableName+" WHERE key = ?", s.getKey(key))
	return err
}

func (s *sqlite) PrefixDelete(ctx context.Context, prefix string) error {
	k := s.getKey(prefix)
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+cacheTableName+" WHERE key LIKE ?", k+"%")
	return err
}

func (s *sqlite) Fetch(ctx context.Context, key string) ([]byte, error) {
	row := s.db.QueryRowContext(ctx, "SELECT value,expire FROM "+cacheTableName+" WHERE key = ?",
		s.getKey(key))
	if row.Err() != nil {
		return nil, row.Err()
	}
	var (
		value   []byte
		expires int64
	)
	err := row.Scan(&value, &expires)
	if err != nil {
		return nil, err
	}
	if expires > 0 && expires < time.Now().Unix() {
		return nil, s.Delete(ctx, key)
	}
	return value, nil
}

func (s *sqlite) FetchMulti(ctx context.Context, keys ...string) ([][]byte, error) {
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
	// 查询
	rows, err := s.db.QueryContext(ctx, "SELECT value,expire FROM "+cacheTableName+" WHERE key IN ("+
		placeholder[:len(placeholder)-1]+
		")", ks...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 扫描结果，过滤过期数据
	i := 0
	ii := make([]int, 0, l)
	for rows.Next() {
		i++
		var (
			value   []byte
			expires int64
		)
		err = rows.Scan(&value, &expires)
		if err != nil {
			return nil, err
		}
		if expires > 0 && expires < time.Now().Unix() {
			ii = append(ii, i-1)
			vals = append(vals, nil)
			continue
		}
		vals = append(vals, value)
	}
	// 删除过期数据
	go func() {
		ks = ks[:0]
		for _, j := range ii {
			ks = append(ks, s.getKey(keys[j]))
		}
		placeholder = strings.Repeat("?,", len(ks))
		_, err = s.db.ExecContext(ctx, "DELETE FROM "+cacheTableName+" WHERE key IN ("+placeholder[:len(placeholder)-1]+")", ks...)
	}()

	return vals, nil
}
func (s *sqlite) Flush(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+cacheTableName)
	return err
}

func (s *sqlite) Save(ctx context.Context, key string, value any, lifeTime time.Duration) error {
	duration := int64(0)
	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}
	var (
		stmt *sql.Stmt
		err  error
	)
	// 查询是否存在
	if s.Contains(ctx, key) {
		stmt, err = s.db.PrepareContext(ctx, "UPDATE "+cacheTableName+" SET value = ?, expire = ? WHERE key = ?")
	} else {
		stmt, err = s.db.PrepareContext(ctx, "INSERT INTO "+cacheTableName+"(value, expire,key) VALUES(?, ?, ?)")
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, value, duration, s.getKey(key))
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
