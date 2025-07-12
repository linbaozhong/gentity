package ace

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
	"golang.org/x/sync/singleflight"
	"sync"

	"github.com/linbaozhong/gentity/pkg/log"
)

type (
	DB struct {
		*sql.DB
		driverName string
		mapper     *reflectx.Mapper

		debug     bool // 如果是调试模式，则打印sql命令及错误
		cacheType cacheType
		cacheOpts any
		cacheMap  sync.Map
		sg        singleflight.Group
	}
)

// Close
func (s *DB) Close() error {
	return s.DB.Close()
}

// Mapper
func (s *DB) Mapper() *reflectx.Mapper {
	return s.mapper
}

// SetDebug
func (s *DB) SetDebug(debug bool) *DB {
	s.debug = debug
	return s
}

// Debug
func (s *DB) Debug() bool {
	return s.debug
}

// SetCache
// opts string: memcache地址(github.com/bradfitz/gomemcache/memcache)
// opts *rd.Options: redis配置(github.com/redis/go-redis/v9)
// opts nil：缺省 sync.Map
func (s *DB) SetCache(t cacheType, opts any) *DB {
	s.cacheType = t
	s.cacheOpts = opts
	return s
}

//
//// Cache
//func (s *DB) Cache(name string) cachego.Cache {
//	if v, ok := s.cacheMap.Load(name); ok {
//		return v.(cachego.Cache)
//	}
//
//	v, _, _ := s.sg.Do(name, func() (any, error) {
//		var v cachego.Cache
//		switch s.cacheType {
//		case CacheTypeMemory:
//			if opts, ok := s.cacheOpts.(string); ok {
//				v = memcached.New(memcache.New(opts), memcached.WithPrefix(name))
//			}
//		case CacheTypeRedis:
//			if opts, ok := s.cacheOpts.(*rd.Options); ok {
//				v = redis.New(rd.NewClient(opts), redis.WithPrefix(name))
//			}
//		default: // CacheTypeSyncMap
//			v = mmap.New() // sync.Map 不需要前缀
//		}
//		if v == nil {
//			v = mmap.New() // sync.Map 不需要前缀
//		}
//		s.cacheMap.Store(name, v)
//		return v, nil
//	})
//	return v.(cachego.Cache)
//}

// Transaction 事务处理
func (s *DB) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	tx, e := s.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}

	var result any
	result, e = f(&Tx{tx, s.mapper, s.Transaction, s.debug})
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
func (s *DB) IsDB() bool {
	return true
}

// QueryContext 执行查询操作
func (s *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext 执行单行查询操作
func (s *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.DB.QueryRowContext(ctx, query, args...)
}

// ExecContext 执行更新、插入、删除等操作
func (s *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(ctx, query, args...)
}

// PrepareContext 为以后的查询或执行创建一个准备好的语句。可以从返回的语句并发地运行多个查询或执行。调用者必须调用语句的Stmt。当不再需要语句时，关闭方法。
// 所提供的上下文用于语句的准备，而不是用于语句的执行。
func (s *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.DB.PrepareContext(ctx, query)
}
