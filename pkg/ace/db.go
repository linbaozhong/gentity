package ace

import (
	"context"
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
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

// Cache
func (s *DB) Cache(name string) cachego.Cache {
	if v, ok := s.cacheMap.Load(name); ok {
		return v.(cachego.Cache)
	}

	v, _, _ := s.sg.Do(name, func() (any, error) {
		var v cachego.Cache
		switch s.cacheType {
		case CacheTypeMemory:
			if opts, ok := s.cacheOpts.(string); ok {
				v = memcached.New(memcache.New(opts), memcached.WithPrefix(name))
			}
		case CacheTypeRedis:
			if opts, ok := s.cacheOpts.(*rd.Options); ok {
				v = redis.New(rd.NewClient(opts), redis.WithPrefix(name))
			}
		default: // CacheTypeSyncMap
			v = mmap.New() // sync.Map 不需要前缀
		}
		if v == nil {
			v = mmap.New() // sync.Map 不需要前缀
		}
		s.cacheMap.Store(name, v)
		return v, nil
	})
	return v.(cachego.Cache)
}

// Transaction 事务处理
func (s *DB) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	tx, e := s.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}

	var result any
	result, e = f(&Tx{tx, s.mapper, s.Cache, s.Transaction, s.debug})
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

//
//func (s *DB) Create() CreateBuilder {
//	return newCreate(s)
//}
//
//func (s *DB) Update() UpdateBuilder {
//	return newUpdate(s)
//}
//
//func (s *DB) Delete() DeleteBuilder {
//	return newDelete(s)
//}
//
//func (s *DB) Select() SelectBuilder {
//	return newSelect(s)
//}
