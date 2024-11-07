package ace

import (
	"context"
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
	"github.com/linbaozhong/gentity/pkg/cachego/memcached"
	"github.com/linbaozhong/gentity/pkg/cachego/redis"
	syc "github.com/linbaozhong/gentity/pkg/cachego/sync"
	"golang.org/x/sync/singleflight"
	"sync"

	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/cachego"
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
		Mapper() *reflectx.Mapper
		// BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Debug() bool
		Cache(string) cachego.Cache
		C(tableName string) *Creator
		D(tableName string) *Deleter
		U(tableName string) *Updater
		R(tableName string) *Selector
	}

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
	Tx struct {
		*sql.Tx
		mapper      *reflectx.Mapper
		cache       func(name string) cachego.Cache
		transaction func(ctx context.Context, f func(tx *Tx) (any, error)) (any, error)
		debug       bool // 如果是调试模式，则打印sql命令及错误
	}
	cacheType string
)

const (
	CacheTypeMemory  cacheType = "memory"
	CacheTypeRedis   cacheType = "redis"
	CacheTypeSyncMap cacheType = "sync"
)

// Connect
func Connect(ctx context.Context, driverName, dns string) (*DB, error) {
	dialect.Register(driverName)
	db, e := sql.Open(driverName, dns)
	if e != nil {
		log.Panic(e)
		return nil, e
	}
	if e = db.Ping(); e != nil {
		db.Close()
		log.Panic(e)
		return nil, e
	}

	obj := &DB{}
	obj.DB = db
	obj.driverName = driverName
	obj.mapper = mapper()
	obj.debug = false

	// 初始化全局的context，使其支持取消
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				obj.Close()
				return
			}
		}
	}()

	return obj, e
}

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
			v = syc.New() // sync.Map 不需要前缀
		}
		if v == nil {
			v = syc.New() // sync.Map 不需要前缀
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

func (t *Tx) Mapper() *reflectx.Mapper {
	return t.mapper
}
func (t *Tx) C(tableName string) *Creator {
	return newCreate(t, tableName)
}
func (t *Tx) U(tableName string) *Updater {
	return NewUpdate(t, tableName)
}
func (t *Tx) D(tableName string) *Deleter {
	return newDelete(t, tableName)
}
func (t *Tx) R(tableName string) *Selector {
	return newSelect(t, tableName)
}
func (t *Tx) Cache(name string) cachego.Cache {
	return t.cache(name)
}
func (t *Tx) Debug() bool {
	return t.debug
}

func (t *Tx) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	return t.transaction(ctx, f)
}
