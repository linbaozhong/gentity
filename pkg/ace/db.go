package ace

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/linbaozhong/gentity/pkg/ace/reflectx"
	"github.com/linbaozhong/gentity/pkg/app"
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
		IsDB() bool
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

var _obj *DB

// Connect
func Connect(driverName, dns string) (*DB, error) {
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

	_obj = &DB{}
	_obj.DB = db
	_obj.driverName = driverName
	_obj.mapper = mapper()
	_obj.debug = false

	app.RegisterServiceCloser(_obj)

	return _obj, e
}

// GetDB
// 调用该方法前，确保已经调用过 Connect 方法并确保没有 error 产生
func GetDB() *DB {
	if _obj == nil {
		log.Panic("db not init")
	}
	return _obj
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
func (s *DB) IsDB() bool {
	return true
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
func (t *Tx) IsDB() bool {
	return false
}
func (t *Tx) Transaction(ctx context.Context, f func(tx *Tx) (any, error)) (any, error) {
	return t.transaction(ctx, f)
}

func Or(fns ...dialect.Condition) dialect.Condition {
	return func() (string, any) {
		if len(fns) == 0 {
			return "", nil
		}
		var (
			buf    bytes.Buffer
			params = make([]any, 0, len(fns))
		)

		buf.WriteString(dialect.Operator_or + "(")
		for i, fn := range fns {
			cond, val := fn()
			if i > 0 {
				buf.WriteString(dialect.Operator_and)
			}
			buf.WriteString(cond)
			if vals, ok := val.([]any); ok {
				params = append(params, vals...)
			} else {
				params = append(params, val)
			}
		}
		buf.WriteString(")")

		return buf.String(), params
	}
}
func And(fns ...dialect.Condition) dialect.Condition {
	return func() (string, any) {
		if len(fns) == 0 {
			return "", nil
		}
		var (
			buf    bytes.Buffer
			params = make([]any, 0, len(fns))
		)

		buf.WriteString(dialect.Operator_and + "(")
		for i, fn := range fns {
			cond, val := fn()
			if i > 0 {
				buf.WriteString(dialect.Operator_or)
			}
			buf.WriteString(cond)
			if vals, ok := val.([]any); ok {
				params = append(params, vals...)
			} else {
				params = append(params, val)
			}
		}
		buf.WriteString(")")

		return buf.String(), params
	}
}

// Order 函数用于创建一个升序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// 并调用 Asc 函数来生成一个升序排序的规则。
// 返回值为一个实现了 dialect.Order 接口的函数，该函数可以被用于指定查询结果的排序方式。
func Order(fs ...dialect.Field) dialect.Order {
	return Asc(fs...)
}

// Asc 函数用于创建一个升序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "ASC" 和指定的字段列表。
// 该函数可用于指定查询结果按指定字段进行升序排序。
func Asc(fs ...dialect.Field) dialect.Order {
	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "ASC" 和字段列表
	return func() (string, []dialect.Field) {
		// 返回升序操作符
		return dialect.Operator_Asc, fs
	}
}

// Desc 函数用于创建一个降序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "DESC" 和指定的字段列表。
// 该函数可用于指定查询结果按指定字段进行降序排序。
func Desc(fs ...dialect.Field) dialect.Order {
	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "DESC" 和字段列表
	return func() (string, []dialect.Field) {
		// 返回降序操作符
		return dialect.Operator_Desc, fs
	}
}
