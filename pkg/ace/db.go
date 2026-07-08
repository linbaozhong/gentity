package ace

import (
	"context"
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
	"github.com/linbaozhong/gentity/pkg/cachego/memcached"
	"github.com/linbaozhong/gentity/pkg/cachego/mmap"
	"github.com/linbaozhong/gentity/pkg/cachego/redis"
	"golang.org/x/sync/singleflight"
	"sync"

	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/log"
	rd "github.com/redis/go-redis/v9"
)

type (
	DB struct {
		*sql.DB
		driverName string
		mapper     *reflectx.Mapper
		dialect    dialect.Dialect
		debug      bool // 如果是调试模式，则打印sql命令及错误
		cacheType  cacheType
		cacheOpts  any
		cacheMap   sync.Map
		sg         singleflight.Group
	}
)

// Close
func (d *DB) Close() error {
	return d.DB.Close()
}

// Mapper
func (d *DB) Mapper() *reflectx.Mapper {
	return d.mapper
}

// Debug 设置调试模式，只打印sql命令，不不执行DAL操作
func (d *DB) Debug(debug ...bool) bool {
	if len(debug) > 0 {
		d.debug = debug[0]
	}
	return d.debug
}

// Dialect
func (d *DB) Dialect() dialect.Dialect {
	return d.dialect
}

// SetCache
// opts string: memcache地址(github.com/bradfitz/gomemcache/memcache)
// opts *rd.Options: redis配置(github.com/redis/go-redis/v9)
// opts nil：缺省 sync.Map
func (d *DB) SetCache(t cacheType, opts any) *DB {
	d.cacheType = t
	d.cacheOpts = opts
	return d
}

// Cache
func (d *DB) Cache(name string) cachego.Cache {
	if v, ok := d.cacheMap.Load(name); ok {
		return v.(cachego.Cache)
	}

	v, _, _ := d.sg.Do(name, func() (any, error) {
		var v cachego.Cache
		switch d.cacheType {
		case CacheTypeMemory:
			if opts, ok := d.cacheOpts.(string); ok {
				v = memcached.New(memcache.New(opts), memcached.WithPrefix(name))
			}
		case CacheTypeRedis:
			if opts, ok := d.cacheOpts.(*rd.Options); ok {
				v = redis.New(rd.NewClient(opts), redis.WithPrefix(name))
			}
		default: // CacheTypeSyncMap
			v = mmap.New() // sync.Map 不需要前缀
		}
		if v == nil {
			v = mmap.New() // sync.Map 不需要前缀
		}
		d.cacheMap.Store(name, v)
		return v, nil
	})
	return v.(cachego.Cache)
}

// Transaction 事务处理
func (d *DB) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	tx, e := d.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}

	var result any
	result, e = f(&Tx{
		tx,
		d.mapper,
		d.Cache,
		d.Transaction,
		d.dialect,
		d.debug},
	)
	if e != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		return result, e
	}

	if e = tx.Commit(); e != nil {
		return result, e
	}

	return result, nil
}
func (d *DB) IsDB() bool {
	return true
}

// QueryContext 执行查询操作
func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext 执行单行查询操作
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRowContext(ctx, query, args...)
}

// ExecContext 执行更新、插入、删除等操作
func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.DB.ExecContext(ctx, query, args...)
}

// PrepareContext 为以后的查询或执行创建一个准备好的语句。可以从返回的语句并发地运行多个查询或执行。调用者必须调用语句的Stmt。当不再需要语句时，关闭方法。
// 所提供的上下文用于语句的准备，而不是用于语句的执行。
func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.DB.PrepareContext(ctx, query)
}
