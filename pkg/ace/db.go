package ace

import (
	"context"
	"database/sql"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/linbaozhong/gentity/pkg/cachego/memcached"
	"github.com/linbaozhong/gentity/pkg/cachego/redis"

	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/cachego"
	syc "github.com/linbaozhong/gentity/pkg/cachego/sync"
	"github.com/linbaozhong/gentity/pkg/log"
	rd "github.com/redis/go-redis/v9"
)

type (
	Cruder interface {
		// C Create 命令体
		C() *Creator
		// R Read 命令体
		R() *Selector
		// U Update 命令体
		U() *Updater
		// D Delete 命令体
		D() *Deleter
	}

	Executer interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Debug() bool
		Cache(prefix string) cachego.Cache
		C(tableName string) *Creator
		D(tableName string) *Deleter
		U(tableName string) *Updater
		R(tableName string) *Selector
	}

	DB struct {
		*sql.DB
		debug     bool // 如果是调试模式，则打印sql命令及错误
		cacheType cacheType
		cacheOpts any
		cacheMap  sync.Map
	}

	cacheType string
	cache     struct {
		prefix string
		cachego.Cache
	}
)

const (
	CacheTypeMemory  cacheType = "memory"
	CacheTypeRedis   cacheType = "redis"
	CacheTypeSyncMap cacheType = "sync"
)

func Connect(driverName, dns string) (*DB, error) {
	dialect.Register(driverName)
	db, e := sql.Open(driverName, dns)
	if e != nil {
		log.Panic(e)
		return nil, e
	}
	if e = db.Ping(); e != nil {
		log.Panic(e)
		return nil, e
	}

	obj := &DB{}
	obj.DB = db
	obj.debug = false

	return obj, e
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
func (s *DB) Cache(prefix string) cachego.Cache {
	if c, ok := s.cacheMap.Load(prefix); ok {
		return c.(*cache).Cache
	}
	var c *cache
	switch s.cacheType {
	case CacheTypeSyncMap:
		c = &cache{prefix, syc.New()}

	case CacheTypeMemory:
		if opts, ok := s.cacheOpts.(string); ok {
			c = &cache{prefix, memcached.New(memcache.New(opts))}
		}
	case CacheTypeRedis:
		if opts, ok := s.cacheOpts.(*rd.Options); ok {
			c = &cache{prefix, redis.New(rd.NewClient(opts))}
		}
	default:
		c = &cache{prefix, syc.New()}
	}
	s.cacheMap.Store(prefix, c)
	return c
}

// Transaction 事务处理
func (s *DB) Transaction(ctx context.Context, f func(tx *sql.Tx) (any, error)) (any, error) {
	tx, e := s.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}

	var result any
	result, e = f(tx)
	if e != nil {
		if e = tx.Rollback(); e != nil {
			log.Error(e)
		}
		return result, e
	}

	if e = tx.Commit(); e != nil {
		return result, e
	}

	return result, nil
}

func (s *DB) C(tableName string) *Creator {
	return newCreate(s, tableName)
}

func (s *DB) U(tableName string) *Updater {
	return NewUpdate(s, tableName)
}

func (s *DB) D(tableName string) *Deleter {
	return newDelete(s, tableName)
}

func (s *DB) R(tableName string) *Selector {
	return newSelect(s, tableName)
}
